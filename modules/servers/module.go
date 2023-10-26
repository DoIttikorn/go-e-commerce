package servers

import (
	"github.com/Doittikorn/go-e-commerce/modules/middlewares/middlewaresHandlers"
	"github.com/Doittikorn/go-e-commerce/modules/middlewares/middlewaresRepositories"
	"github.com/Doittikorn/go-e-commerce/modules/middlewares/middlewaresUsecases"
	"github.com/Doittikorn/go-e-commerce/modules/monitor/monitorHandlers"
	"github.com/gofiber/fiber/v2"
)

type ModuleFactoryImpl interface {
	MonitorModule()
	UsersModule()
	AppinfoModule()
	FilesModule() IFilesModule
	ProductsModule() IProductsModule
	OrdersModule()
}

type moduleFactory struct {
	router fiber.Router
	server *server
	mid    middlewaresHandlers.MiddlewaresHandlerImpl
}

func InitModule(r fiber.Router, s *server, m middlewaresHandlers.MiddlewaresHandlerImpl) ModuleFactoryImpl {
	return &moduleFactory{
		router: r,
		server: s,
		mid:    m,
	}
}

func InitMiddlewares(s *server) middlewaresHandlers.MiddlewaresHandlerImpl {
	repository := middlewaresRepositories.MiddlewaresRepositry(s.db)
	usecase := middlewaresUsecases.MiddlewaresUsecase(repository)
	return middlewaresHandlers.MiddlewaresHandler(s.cfg, usecase)
}

func (m *moduleFactory) MonitorModule() {
	monitor := monitorHandlers.New(m.server.cfg)
	m.router.Get("/health-check", monitor.HealthCheck)
}
