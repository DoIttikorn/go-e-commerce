package servers

import (
	"github.com/Doittikorn/go-e-commerce/modules/orders/ordersHandlers"
	"github.com/Doittikorn/go-e-commerce/modules/orders/ordersRepositories"
	"github.com/Doittikorn/go-e-commerce/modules/orders/ordersUsecases"
)

func (m *moduleFactory) OrdersModule() {
	ordersRepository := ordersRepositories.OrdersRepository(m.server.db)
	ordersUsecase := ordersUsecases.OrdersUsecase(ordersRepository, m.ProductsModule().Repository())
	ordersHandler := ordersHandlers.OrdersHandler(m.server.cfg, ordersUsecase)

	router := m.router.Group("/orders")

	router.Post("/", m.mid.JwtAuth(), ordersHandler.InsertOrder)

	router.Get("/", m.mid.JwtAuth(), m.mid.Authorize(2), ordersHandler.FindOrder)
	router.Get("/:user_id/:order_id", m.mid.JwtAuth(), m.mid.VerifyParamUserId(), ordersHandler.FindOneOrder)

	router.Patch("/:user_id/:order_id", m.mid.JwtAuth(), m.mid.VerifyParamUserId(), ordersHandler.UpdateOrder)
}
