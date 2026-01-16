package router

import (
	"net/http"

	"github.com/a5932016/go-ddd-example/util/mGin"
	"github.com/gin-gonic/gin"
)

type bindIdURI struct {
	ID uint `uri:"id" binding:"required,number"`
}

func (rH Handler) getUserHandler(c *gin.Context) {
	ctx := mGin.NewContext(c)
	var boundIdURI bindIdURI

	if err := ctx.ShouldBindUri(&boundIdURI); err != nil {
		ctx.WithError(err).Response(http.StatusBadRequest, "Invalid URI")
		return
	}

	user, err := rH.handler.GetUser(ctx, boundIdURI.ID)
	if err != nil {
		ctx.WithError(err).Response(http.StatusInternalServerError, "rH.handler.GetUser")
		return
	}

	ctx.WithData(user).Response(http.StatusOK, "")
	return
}
