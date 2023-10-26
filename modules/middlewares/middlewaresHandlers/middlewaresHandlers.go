package middlewaresHandlers

import (
	"net/http"
	"strings"

	"github.com/Doittikorn/go-e-commerce/config"
	"github.com/Doittikorn/go-e-commerce/modules/entities"
	"github.com/Doittikorn/go-e-commerce/modules/middlewares/middlewaresUsecases"
	"github.com/Doittikorn/go-e-commerce/pkg/auth"
	"github.com/Doittikorn/go-e-commerce/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

type middlewareHandlersErrCode string

const (
	routerCheckErr middlewareHandlersErrCode = "middleware-001"
	jwtAuthErr     middlewareHandlersErrCode = "middleware-002"
	paramsCheckErr middlewareHandlersErrCode = "middleware-003"
	authorizeErr   middlewareHandlersErrCode = "middleware-004"
	apiKeyErr      middlewareHandlersErrCode = "middlware-005"
)

type MiddlewaresHandlerImpl interface {
	Cors() fiber.Handler
	RouterCheck() fiber.Handler
	Logger() fiber.Handler
	JwtAuth() fiber.Handler
	VerifyParamUserId() fiber.Handler
	Authorize(...int) fiber.Handler
	ApiKeyAuth() fiber.Handler
	StreamingFile() fiber.Handler
}

type middlewaresHandler struct {
	cfg               config.ConfigImpl
	middlewareUsecase middlewaresUsecases.MiddlewaresUsecaseImpl
}

func MiddlewaresHandler(cfg config.ConfigImpl, middlewareU middlewaresUsecases.MiddlewaresUsecaseImpl) MiddlewaresHandlerImpl {

	return &middlewaresHandler{
		middlewareUsecase: middlewareU,
		cfg:               cfg,
	}
}

func (h *middlewaresHandler) Cors() fiber.Handler {
	return cors.New(cors.Config{
		Next:             cors.ConfigDefault.Next,
		AllowOrigins:     "*",
		AllowHeaders:     "GET,POST,HEAD,PUT,DELETE,PATCH",
		AllowCredentials: false,
		ExposeHeaders:    "",
		MaxAge:           0,
	})
}

// ใช้ในการตรวจสอบว่า router ที่เรียกมานั้นมีอยู่หรือไม่
func (h *middlewaresHandler) RouterCheck() fiber.Handler {

	return func(c *fiber.Ctx) error {
		return entities.NewResponse(c).Error(
			http.StatusNotFound,
			string(routerCheckErr),
			"router not found",
		).Res()
	}
}

// ใช้ในการเก็บ log ของ request ที่เข้ามา
func (h *middlewaresHandler) Logger() fiber.Handler {
	return logger.New(logger.Config{
		Format:     "${time} [${ip}] ${status} - ${method} ${path}\n",
		TimeFormat: "02/01/2006",
		TimeZone:   "Asia/Bangkok",
	})
}

// ตรวจสอบ token ว่าถูกสร้างขึ้นโดยเราหรือไม่ ที่ oauth table
func (h *middlewaresHandler) JwtAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := strings.TrimPrefix(c.Get("Authorization"), "Bearer ")

		result, err := auth.ParseToken(h.cfg.JWT(), token)

		if err != nil {
			return entities.NewResponse(c).Error(
				http.StatusUnauthorized,
				string(jwtAuthErr),
				"unauthorized",
			).Res()
		}

		claims := result.Claims

		if !h.middlewareUsecase.FindAccessToken(claims.Id, token) {
			return entities.NewResponse(c).Error(
				http.StatusUnauthorized,
				string(jwtAuthErr),
				"no permission",
			).Res()

		}

		// Set information user to context
		c.Locals("userId", claims.Id)
		c.Locals("userRoleId", claims.RoleId)

		return c.Next()
	}
}

// ตรวจสอบว่า path parameter userId ที่ส่งมาตรงกับ userId ที่อยู่ใน token หรือไม่
func (h *middlewaresHandler) VerifyParamUserId() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userId, ok := c.Locals("userId").(string)
		if !ok {
			return entities.NewResponse(c).Error(
				http.StatusUnauthorized,
				string(paramsCheckErr),
				"unauthorized",
			).Res()
		}
		if c.Params("userId") != userId {
			return entities.NewResponse(c).Error(
				http.StatusUnauthorized,
				string(paramsCheckErr),
				"permission denied",
			).Res()
		}
		return c.Next()
	}
}

// ตรวจสอบว่า user ที่ส่งมามีสิทธิ์ในการเข้าถึงหรือไม่
func (h *middlewaresHandler) Authorize(expectedRoleId ...int) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRoleId, ok := c.Locals("userRoleId").(int)
		if !ok {
			return entities.NewResponse(c).Error(
				http.StatusUnauthorized,
				string(authorizeErr),
				"userId is invalid",
			).Res()
		}
		roles, err := h.middlewareUsecase.FindRole()
		if err != nil {
			return entities.NewResponse(c).Error(
				http.StatusInternalServerError,
				string(authorizeErr),
				err.Error(),
			).Res()
		}

		sum := 0
		for _, v := range expectedRoleId {
			sum += v
		}

		expectedValueBinary := utils.BinaryConverter(sum, len(roles))

		userValueBinary := utils.BinaryConverter(userRoleId, len(roles))

		// user -> 0 1 0
		// expected -> 1 1 0

		for i := range userValueBinary {
			if userValueBinary[i]&expectedValueBinary[i] == 1 {
				return c.Next()
			}
		}
		return entities.NewResponse(c).Error(
			fiber.ErrUnauthorized.Code,
			string(authorizeErr),
			"no permission to access",
		).Res()

	}
}

func (h *middlewaresHandler) ApiKeyAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		key := c.Get("X-Api-Key")
		if _, err := auth.ParseApiKey(h.cfg.JWT(), key); err != nil {
			return entities.NewResponse(c).Error(
				fiber.ErrUnauthorized.Code,
				string(apiKeyErr),
				"apikey is invalid or required",
			).Res()
		}
		return c.Next()
	}
}

// Streaming file
func (h *middlewaresHandler) StreamingFile() fiber.Handler {
	return filesystem.New(filesystem.Config{
		Root: http.Dir("./assets/images"),
	})
}
