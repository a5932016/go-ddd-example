package router

import (
	"net/http"

	"github.com/a5932016/go-ddd-example/model"
	"github.com/a5932016/go-ddd-example/util/mGin"
	"github.com/gin-gonic/gin"
)

type appRouter struct {
	method        string
	endpoint      string
	allowancePair allowancePair
	worker        gin.HandlerFunc
}

type allowancePair struct {
	Resource            model.Resource
	Action              model.Action
	RootOnly            bool
	SelfPrivilege       bool // param ID required
	SelfInterdictFilter bool // param ID required
	HierarchyFilter     bool // param ID required
}

func (rH Handler) getRouter() (routes []appRouter) {
	return []appRouter{
		// user
		appRouter{http.MethodGet, "/user/:id", allowancePair{Resource: model.ResourceUser, Action: model.ActionRead, SelfPrivilege: true}, rH.getUserHandler},
	}
}

func (rH Handler) healthHandler(c *gin.Context) {
	ctx := mGin.NewContext(c)
	ctx.WithData(struct{}{}).Response(http.StatusOK, "")
}
