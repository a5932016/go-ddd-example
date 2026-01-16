package migration

import (
	"github.com/a5932016/go-ddd-example/model"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

var firstMigration = &gormigrate.Migration{
	ID: "firstMigration",
	Migrate: func(db *gorm.DB) error {
		return db.AutoMigrate(&model.User{})
	},
	Rollback: func(db *gorm.DB) error {
		return nil
	},
}
