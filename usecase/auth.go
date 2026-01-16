package usecase

import (
	"context"
	"fmt"

	"github.com/a5932016/go-ddd-example/customerror"
	"github.com/a5932016/go-ddd-example/model"
	"github.com/a5932016/go-ddd-example/singleton/session"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func (h HandlerConstructor) Login(c context.Context, account, password string) (sessionId string, user model.User, err error) {
	// Get user
	user, err = h.dbRepo.GetUserByAccount(account)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", model.User{}, customerror.AccountNotFound
		}
		return "", model.User{}, errors.Wrap(err, "dbRepo.GetUserByAccount")
	}

	// Verify password
	if !model.VerifyPassword(password, user.Password) {
		return "", model.User{}, customerror.WrongPassword
	}

	tx := h.dbRepo.Begin()
	defer tx.Rollback()

	// txPer, closeTx, err := h.perRepo.BeginWithTx(tx.DB())
	// if err != nil {
	// 	return "", model.User{}, errors.Wrap(err, "perRepo.BeginWithTx")
	// }
	// defer closeTx(c)

	// for index := range user.Divisions {
	// 	policies, err := txPer.GetPolicies(user.Divisions[index].GetPrefixedNameID())
	// 	if err != nil {
	// 		return "", model.User{}, errors.Wrap(err, "txPer.GetPolicies")
	// 	}
	// 	user.Divisions[index].Permissions, err = h.permissionsHandler.CasbinPoliciesToPermissions(policies)
	// 	if err != nil {
	// 		return "", model.User{}, errors.Wrap(err, "model.CasbinPoliciesToPermissions")
	// 	}
	// }

	// Set session
	sc, err := h.sessionManager.SessionStart("")
	if err != nil {
		return "", model.User{}, err
	}

	userStr, err := h.stringifyUser(user)
	if err != nil {
		return "", model.User{}, errors.Wrap(err, "stringifyUser")
	}

	if err := sc.Session.Set(SIDUser, userStr); err != nil {
		return "", model.User{}, errors.Wrap(err, "Session.Set(user)")
	}

	return sc.Session.SessionID(), user, nil
}

const (
	SID     = "sid"
	SIDUser = "user"
)

func (h HandlerConstructor) Logout(c context.Context) (err error) {
	sid := c.Value(SID).(string)
	err = h.sessionManager.SessionDestroy(sid)
	return
}

func (h HandlerConstructor) ForgotPassword(c *gin.Context, account string) (resetToken string, err error) {
	// Get approver
	approver, err := h.getRequestUser(c)
	if err != nil {
		return "", err
	}

	// Check hierarchy permission
	if !approver.IsRoot {
		return "", customerror.NoHierarchyPermission
	}

	// Get user
	user, err := h.dbRepo.GetUserByAccount(account)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", customerror.AccountNotFound
		}
		return "", errors.Wrap(err, "dbRepo.GetUserByAccount")
	}

	// Set session
	sc, err := h.sessionManager.SessionStart("")
	if err != nil {
		return "", err
	}

	userStr, err := h.stringifyUser(user)
	if err != nil {
		return "", errors.Wrap(err, "stringifyUser")
	}

	if err := sc.Session.Set(SIDUser, userStr); err != nil {
		return "", errors.Wrap(err, "Session.Set(user)")
	}

	approverStr, err := h.stringifyUser(approver)
	if err != nil {
		return "", errors.Wrap(err, "stringifyUser")
	}

	if err := sc.Session.Set("resetPwdApprover", approverStr); err != nil {
		return "", errors.Wrap(err, "Session.Set(user)")
	}

	if err := sc.Session.Set("resetPwdApproverIP", c.ClientIP()); err != nil {
		return "", errors.Wrap(err, "Session.Set(user)")
	}

	return sc.Session.SessionID(), nil
}

func (h HandlerConstructor) ResetPassword(c context.Context, resetToken, password string) error {
	user, err := h.GetRequestUserFromSID(resetToken)
	if err != nil {
		return err
	}

	// Change password
	user.Password, err = h.hashPassword(password)
	if err != nil {
		return err
	}

	if err := h.dbRepo.UpdateUserPassword(user.ID, user.Password); err != nil {
		return err
	}

	return h.sessionManager.SessionDestroy(resetToken)
}

func (h HandlerConstructor) stringifyUser(user model.User) (string, error) {
	userBytes, err := json.Marshal(user)
	if err != nil {
		return "", errors.Wrap(err, "json.Marshal(user)")
	}

	return string(userBytes), nil
}

func (h HandlerConstructor) getRequestUser(c context.Context) (model.User, error) {
	sid := c.Value(SID).(string)
	if !(len(sid) > 0) {
		return model.User{}, errors.New("sid not found")
	}

	return h.GetRequestUserFromSID(sid)
}

func (h HandlerConstructor) GetRequestUserFromSID(sessionID string) (model.User, error) {
	sc, err := h.sessionManager.SessionStart(sessionID)
	if err != nil {
		if errors.Is(err, session.ErrSessionNotExisted) {
			return model.User{}, customerror.InvalidSession
		}
		return model.User{}, errors.Wrap(err, "sessionManager.SessionStart")
	}

	userStr, ok := sc.Session.Get(SIDUser).(string)
	if !ok {
		return model.User{}, errors.Wrap(err, "Session.Get(user).(string)")
	}

	var requester model.User
	if err := json.Unmarshal([]byte(userStr), &requester); err != nil {
		return model.User{}, errors.Wrap(err, "json.Unmarshal(user)")
	}

	return requester, nil
}

func (h HandlerConstructor) hashPassword(password string) (string, error) {
	passwordBytes, err := model.HashPassword(password)
	if err != nil {
		if err == bcrypt.ErrPasswordTooLong {
			return "", customerror.PasswordTooLong
		}

		return "", errors.Wrap(err, fmt.Sprintf("model.HashPassword(%s)", password))
	}

	return string(passwordBytes), nil
}
