//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/afero"
	"gorm.io/gorm"

	"github.com/a5932016/go-ddd-example/migration"
	"github.com/a5932016/go-ddd-example/router"
)

// InitMigration init router
func InitMigration(mySqlC *gorm.DB) (migration.Migration, error) {
	wire.Build(
		entityHandler,
		dbRepoProvider,
		perRepoProvider,
		migration.New,
	)
	return migration.Migration{}, nil
}

// InitRouter init router
func InitRouter(appFs afero.Fs, mySqlC *gorm.DB, redisC *redis.Client) (router.Handler, error) {
	wire.Build(
		entityHandler,
		entityUseCaseHandler,
		repositoryProvider,
		usecaseProvider,
		permissionsHandler,
		sessionRedisProviderManager,
		router.NewRouter,
	)
	return router.Handler{}, nil
}
