package usecase

import (
	"context"

	"github.com/a5932016/go-ddd-example/customerror"
	"github.com/a5932016/go-ddd-example/model"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

func (h HandlerConstructor) GetUser(c context.Context, id uint) (user model.User, err error) {
	user, err = h.dbRepo.GetUser(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.User{}, customerror.RecordNotFound
		}
		return model.User{}, errors.Wrap(err, "dbRepo.GetUser")
	}
	return
}
