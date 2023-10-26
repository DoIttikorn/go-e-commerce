package usersHandlers

import (
	"net/http"
	"strings"

	"github.com/Doittikorn/go-e-commerce/config"
	"github.com/Doittikorn/go-e-commerce/modules/entities"
	"github.com/Doittikorn/go-e-commerce/modules/users"
	"github.com/Doittikorn/go-e-commerce/modules/users/usersUsecases"
	"github.com/Doittikorn/go-e-commerce/pkg/auth"
	"github.com/gofiber/fiber/v2"
)

type userHandlerErrcode string

const (
	signUpCustomerErr     userHandlerErrcode = "users_handler_001"
	signInErr             userHandlerErrcode = "users_handler_002"
	refreshPasportErr     userHandlerErrcode = "users_handler_003"
	signOutErr            userHandlerErrcode = "users_handler_004"
	signUpAdminErr        userHandlerErrcode = "users_handler_005"
	generateAdminTokenErr userHandlerErrcode = "users_handler_006"
	getUserProfileErr     userHandlerErrcode = "users_handler_007"
)

type UsersHandlersImpl interface {
	SignUpCustomer(c *fiber.Ctx) error
	SignIn(c *fiber.Ctx) error
	RefreshPasport(c *fiber.Ctx) error
	SignOut(c *fiber.Ctx) error
	SignUpAdmin(c *fiber.Ctx) error
	GenerateAdminToken(c *fiber.Ctx) error
	GetUserProfile(c *fiber.Ctx) error
}

type usersHandler struct {
	cfg          config.ConfigImpl
	usersUsecase usersUsecases.UsersUsecasesImpl
}

func New(cfg config.ConfigImpl, userUsecase usersUsecases.UsersUsecasesImpl) UsersHandlersImpl {
	return &usersHandler{
		cfg:          cfg,
		usersUsecase: userUsecase,
	}
}

func (h *usersHandler) SignUpCustomer(c *fiber.Ctx) error {
	// Requset body parser
	req := new(users.UserRegisterReq)
	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(http.StatusBadRequest, string(signUpCustomerErr), err.Error()).Res()
	}

	// Email validation
	if !req.IsEmail() {
		return entities.NewResponse(c).Error(
			http.StatusBadRequest,
			string(signUpCustomerErr),
			"invalid email",
		).Res()
	}

	result, err := h.usersUsecase.InsertCustomer(req)
	if err != nil {
		switch err.Error() {
		case "username has been used":
			return entities.NewResponse(c).Error(http.StatusBadRequest, string(signUpCustomerErr), err.Error()).Res()
		case "email has been used":
			return entities.NewResponse(c).Error(http.StatusBadRequest, string(signUpCustomerErr), err.Error()).Res()
		default:
			return entities.NewResponse(c).Error(http.StatusInternalServerError, string(signUpCustomerErr), err.Error()).Res()

		}
	}
	return entities.NewResponse(c).Success(http.StatusCreated, result).Res()
}

func (h *usersHandler) SignUpAdmin(c *fiber.Ctx) error {
	// Requset body parsers
	req := new(users.UserRegisterReq)
	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(http.StatusBadRequest, string(signUpCustomerErr), err.Error()).Res()
	}

	// Email validation
	if !req.IsEmail() {
		return entities.NewResponse(c).Error(
			http.StatusBadRequest,
			string(signUpCustomerErr),
			"invalid email",
		).Res()
	}

	// Insert admin
	result, err := h.usersUsecase.InsertAdmin(req)
	if err != nil {
		switch err.Error() {
		case "username has been used":
			return entities.NewResponse(c).Error(http.StatusBadRequest, string(signUpCustomerErr), err.Error()).Res()
		case "email has been used":
			return entities.NewResponse(c).Error(http.StatusBadRequest, string(signUpCustomerErr), err.Error()).Res()
		default:
			return entities.NewResponse(c).Error(http.StatusInternalServerError, string(signUpCustomerErr), err.Error()).Res()

		}
	}
	return entities.NewResponse(c).Success(http.StatusCreated, result).Res()
}

func (h *usersHandler) SignIn(c *fiber.Ctx) error {
	req := new(users.UserCredential)
	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(
			http.StatusBadRequest,
			string(signInErr),
			err.Error(),
		).Res()
	}
	passport, err := h.usersUsecase.GetPassport(req)
	if err != nil {
		return entities.NewResponse(c).Error(
			http.StatusBadRequest,
			string(signInErr),
			err.Error(),
		).Res()
	}
	return entities.NewResponse(c).Success(http.StatusOK, passport).Res()
}

func (h *usersHandler) RefreshPasport(c *fiber.Ctx) error {
	req := new(users.UserRefreshCredential)
	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(
			http.StatusBadRequest,
			string(signInErr),
			err.Error(),
		).Res()
	}
	passport, err := h.usersUsecase.RefreshPassport(req)
	if err != nil {
		return entities.NewResponse(c).Error(
			http.StatusBadRequest,
			string(refreshPasportErr),
			err.Error(),
		).Res()
	}
	return entities.NewResponse(c).Success(http.StatusOK, passport).Res()
}

func (h *usersHandler) SignOut(c *fiber.Ctx) error {
	req := new(users.UserRemoveCredential)
	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(
			http.StatusBadRequest,
			string(signOutErr),
			err.Error(),
		).Res()
	}

	if err := h.usersUsecase.DeleteOauth(req.OathId); err != nil {
		return entities.NewResponse(c).Error(
			http.StatusBadRequest,
			string(signOutErr),
			err.Error(),
		).Res()
	}
	return entities.NewResponse(c).Success(http.StatusOK, nil).Res()
}

func (h *usersHandler) GenerateAdminToken(c *fiber.Ctx) error {
	adminToken, err := auth.New(auth.Admin, h.cfg.JWT(), nil)
	if err != nil {
		return entities.NewResponse(c).Error(
			http.StatusInternalServerError,
			string(generateAdminTokenErr),
			err.Error(),
		).Res()
	}
	return entities.NewResponse(c).Success(http.StatusOK, &struct {
		Token string `json:"token"`
	}{
		Token: adminToken.SignToken(),
	}).Res()
}

func (h *usersHandler) GetUserProfile(c *fiber.Ctx) error {
	// get params
	userId := strings.Trim(c.Params("userId"), " ")
	result, err := h.usersUsecase.GetUserProfile(userId)
	if err != nil {
		switch err.Error() {
		case "get user failed: sql: no rows in result set":
			return entities.NewResponse(c).Error(http.StatusBadRequest, string(getUserProfileErr), "resource not found").Res()
		default:
			return entities.NewResponse(c).Error(http.StatusInternalServerError, string(getUserProfileErr), err.Error()).Res()
		}
	}
	return entities.NewResponse(c).Success(http.StatusOK, result).Res()
}
