package usecase

import (
	"context"

	"github.com/a5932016/go-ddd-example/model"
	"github.com/a5932016/go-ddd-example/repository"
	"github.com/a5932016/go-ddd-example/repository/casbin"
	"github.com/a5932016/go-ddd-example/repository/fs"
	"github.com/a5932016/go-ddd-example/singleton/session"
	"github.com/gin-gonic/gin"
)

// Handler handler
type Handler interface {
	Auth
	User
}

// NewHandler new handler
func NewHandler(
	dbRepo repository.DBRepository,
	memRepo repository.MemRepository,
	perRepo *casbin.PERRepository,
	fsRepo fs.FSRepository,
	sessionManager *session.Manager,
	permissionsHandler model.PermissionsHandler,
) *HandlerConstructor {
	h := &HandlerConstructor{
		dbRepo:             dbRepo,
		memRepo:            memRepo,
		perRepo:            perRepo,
		fsRepo:             fsRepo,
		sessionManager:     sessionManager,
		permissionsHandler: permissionsHandler,
	}

	return h
}

// HandlerConstructor HandlerConstructor
type HandlerConstructor struct {
	dbRepo             repository.DBRepository
	memRepo            repository.MemRepository
	perRepo            *casbin.PERRepository
	fsRepo             fs.FSRepository
	sessionManager     *session.Manager
	permissionsHandler model.PermissionsHandler
}

type Auth interface {
	Login(c context.Context, account, password string) (sessionId string, user model.User, err error)
	Logout(c context.Context) (err error)
	ForgotPassword(c *gin.Context, account string) (resetToken string, err error)
	ResetPassword(c context.Context, resetToken, password string) error
}

type User interface {
	GetUser(c context.Context, id uint) (user model.User, err error)
	GetRequestUserFromSID(sessionID string) (model.User, error)
}
