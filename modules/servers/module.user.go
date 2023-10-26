package servers

import (
	"github.com/Doittikorn/go-e-commerce/modules/users/usersHandlers"
	"github.com/Doittikorn/go-e-commerce/modules/users/usersRepositories"
	"github.com/Doittikorn/go-e-commerce/modules/users/usersUsecases"
)

func (m *moduleFactory) UsersModule() {
	repository := usersRepositories.New(m.server.db)
	usecase := usersUsecases.New(m.server.cfg, repository)
	handler := usersHandlers.New(m.server.cfg, usecase)

	router := m.router.Group("/users")
	router.Post("/signin", handler.SignIn)
	router.Post("/signup", handler.SignUpCustomer)
	router.Post("/refresh", handler.RefreshPasport)
	router.Delete("/signout", handler.SignOut)
	router.Post("/signup-admin", handler.SignUpAdmin)

	router.Get("/:userId", m.mid.JwtAuth(), m.mid.VerifyParamUserId(), handler.GetUserProfile)
	router.Get("/admin/secret", m.mid.JwtAuth(), m.mid.Authorize(2), handler.GenerateAdminToken)
}
