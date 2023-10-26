package monitorHandlers

import (
	"net/http"

	"github.com/Doittikorn/go-e-commerce/config"
	"github.com/Doittikorn/go-e-commerce/modules/entities"
	"github.com/Doittikorn/go-e-commerce/modules/monitor"
	"github.com/gofiber/fiber/v2"
)

type MonitorHandlersImpl interface {
	HealthCheck(c *fiber.Ctx) error
}

type monitorHandlers struct {
	cfg config.ConfigImpl
}

func New(cfg config.ConfigImpl) MonitorHandlersImpl {
	return &monitorHandlers{
		cfg: cfg,
	}
}

func (m *monitorHandlers) HealthCheck(c *fiber.Ctx) error {
	res := &monitor.Monitor{
		Name:    m.cfg.App().Name(),
		Version: m.cfg.App().Version(),
	}
	return entities.NewResponse(c).Success(http.StatusOK, res).Res()
}
