package middlewaresUsecases

import (
	"github.com/Doittikorn/go-e-commerce/modules/middlewares"
	"github.com/Doittikorn/go-e-commerce/modules/middlewares/middlewaresRepositories"
)

type MiddlewaresUsecaseImpl interface {
	FindAccessToken(userId, accessToken string) bool
	FindRole() ([]*middlewares.Role, error)
}

type middlewaresUsecases struct {
	middlewaresRepository middlewaresRepositories.MiddlewaresRepositoryImpl
}

func MiddlewaresUsecase(middlewareRepository middlewaresRepositories.MiddlewaresRepositoryImpl) MiddlewaresUsecaseImpl {

	return &middlewaresUsecases{
		middlewaresRepository: middlewareRepository,
	}
}

func (u *middlewaresUsecases) FindAccessToken(userId, accessToken string) bool {
	return u.middlewaresRepository.FindAccessToken(userId, accessToken)
}

func (u *middlewaresUsecases) FindRole() ([]*middlewares.Role, error) {
	return u.middlewaresRepository.FindRole()
}
