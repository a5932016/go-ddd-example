package repository

import (
	"github.com/a5932016/go-ddd-example/model"
	"github.com/a5932016/go-ddd-example/singleton/entity"

	"gorm.io/gorm"
)

// DBRepository interface
type DBRepository interface {
	RDBMS
	User
}

type RDBMS interface {
	DB() *gorm.DB
	Migrate(fn func(*gorm.DB) error)
	Debug()
	EntityCtrl() *entity.EntityHandler

	// transaction
	Begin() DBRepository
	Commit() error
	Rollback() error
}

type User interface {
	GetUser(id uint) (user model.User, err error)
	GetUserByAccount(email string) (user model.User, err error)
	UpdateUserPassword(id uint, password string) error
}
