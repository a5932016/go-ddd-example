package router

import (
	"net/http"
	"strings"

	"github.com/a5932016/go-ddd-example/customerror"
	"github.com/a5932016/go-ddd-example/model"
	"github.com/a5932016/go-ddd-example/util/log"
	"github.com/a5932016/go-ddd-example/util/mGin"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
)

type loginBody struct {
	Account  string `json:"account" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (rH Handler) loginHandler(c *gin.Context) {
	ctx := mGin.NewContext(c)

	var body loginBody
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.WithError(err).Response(http.StatusBadRequest, "Invalid JSON")
		return
	}

	sessionId, user, err := rH.handler.Login(ctx, body.Account, body.Password)
	if err != nil {
		ctx.WithError(err).Response(http.StatusInternalServerError, "handler.Login")
		return
	}

	ctx.WithData(map[string]interface{}{
		"sessionId": sessionId,
		"user":      user,
	}).Response(http.StatusOK, "")

	return
}

func (rH Handler) logoutHandler(c *gin.Context) {
	ctx := mGin.NewContext(c)

	if err := rH.handler.Logout(ctx); err != nil {
		log.FromContext(c).WithError(err).Error("handler.Logout")
	}

	ctx.WithData(struct{}{}).Response(http.StatusOK, "")
	return
}

type forgotPasswordBody struct {
	Email string `json:"email" binding:"required"`
}

func (rH Handler) forgotPasswordHandler(c *gin.Context) {
	ctx := mGin.NewContext(c)

	var body forgotPasswordBody
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.WithError(err).Response(http.StatusBadRequest, "Invalid JSON")
		return
	}

	if ok := govalidator.IsEmail(body.Email); !ok {
		ctx.ResponseWithCustomError(customerror.InvalidEmail)
		return
	}

	resetToken, err := rH.handler.ForgotPassword(c, body.Email)
	if err != nil {
		ctx.WithError(err).Response(http.StatusInternalServerError, "handler.ForgotPassword")
		return
	}

	ctx.WithData(map[string]interface{}{
		"resetToken": resetToken,
	}).Response(http.StatusOK, "")
	return
}

type resetPasswordQuery struct {
	ResetToken string `form:"resetToken" binding:"required"`
}

type resetPasswordBody struct {
	Password             string `json:"password" binding:"required"`
	PasswordConfirmation string `json:"passwordConfirmation" binding:"required"`
}

func (rH Handler) resetPasswordHandler(c *gin.Context) {
	ctx := mGin.NewContext(c)

	var query resetPasswordQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		ctx.WithError(err).Response(http.StatusBadRequest, "Invalid Query")
		return
	}
	var body resetPasswordBody
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.WithError(err).Response(http.StatusBadRequest, "Invalid JSON")
		return
	}

	// Password
	if err := model.ValidatePassword(body.Password); err != nil {
		copyCustomErr := customerror.InvalidPassword
		copyCustomErr.Message = err.Error()
		ctx.ResponseWithCustomError(copyCustomErr)
		return
	}
	if !strings.EqualFold(body.Password, body.PasswordConfirmation) {
		ctx.ResponseWithCustomError(customerror.InvalidPasswordConfirmation)
		return
	}

	if err := rH.handler.ResetPassword(ctx, query.ResetToken, body.Password); err != nil {
		ctx.WithError(err).Response(http.StatusInternalServerError, "handler.ResetPassword")
		return
	}

	ctx.WithData(struct{}{}).Response(http.StatusOK, "")
	return
}
