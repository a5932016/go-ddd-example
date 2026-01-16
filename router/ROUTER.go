package router

import (
	"github.com/a5932016/go-ddd-example/repository"
	"github.com/a5932016/go-ddd-example/repository/casbin"
	"github.com/a5932016/go-ddd-example/singleton/entityUsecase"
	"github.com/a5932016/go-ddd-example/singleton/session"
	"github.com/a5932016/go-ddd-example/usecase"
)

// Handler router handler
type Handler struct {
	handler        usecase.Handler
	entityHandler  entityUsecase.EntityUseCase
	memRepo        repository.MemRepository
	perRepo        *casbin.PERRepository
	sessionManager *session.Manager
}

// NewRouter new router handler
func NewRouter(
	handler usecase.Handler,
	entityHandler entityUsecase.EntityUseCase,
	memRepo repository.MemRepository,
	perRepo *casbin.PERRepository,
	sessionManager *session.Manager,
) Handler {
	return Handler{
		handler:        handler,
		entityHandler:  entityHandler,
		memRepo:        memRepo,
		perRepo:        perRepo,
		sessionManager: sessionManager,
	}
}
