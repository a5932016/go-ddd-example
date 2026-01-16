package mysql

import (
	"github.com/a5932016/go-ddd-example/model"
	"github.com/a5932016/go-ddd-example/util/mGorm"
	"github.com/pkg/errors"
)

func (s *DBRepository) GetUser(id uint) (user model.User, err error) {
	mDB := mGorm.New(s.db.Model(&model.User{}))
	mDB.DB = mDB.DB.Where("id = ?", id)
	err = mDB.DB.Find(&user).Error
	return
}

func (s *DBRepository) GetUserByAccount(email string) (user model.User, err error) {
	if err = s.db.
		Where("email = ?", email).
		First(&user).Error; err != nil {
		err = errors.Wrap(err, "Failed to select user")
		return
	}
	return
}

func (s *DBRepository) UpdateUserPassword(id uint, password string) error {
	return s.db.Model(&model.User{}).Where("id = ?", id).Update("password", password).Error
}
