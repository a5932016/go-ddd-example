package main

import (
	"github.com/a5932016/go-ddd-example/config"
	"github.com/a5932016/go-ddd-example/model"
	"github.com/a5932016/go-ddd-example/singleton/entity"
	"github.com/a5932016/go-ddd-example/singleton/entity/eGorm"
	"github.com/a5932016/go-ddd-example/singleton/entityUsecase"
	"github.com/a5932016/go-ddd-example/singleton/session"
	redisProvider "github.com/a5932016/go-ddd-example/singleton/session/provider/redis"
	"github.com/google/wire"
	"gorm.io/gorm"
)

func _sessionRedisProviderSessionName() session.SessionName {
	return session.SessionName(config.Env.SessionAuth.Name)
}

func _sessionRedisProviderMaxLifeTime() session.MaxLifeTime {
	return session.MaxLifeTime(config.Env.SessionAuth.MaxLifeTime)
}

var (
	sessionRedisProviderManager = wire.NewSet(
		_sessionRedisProvider,
		_sessionRedisProviderSessionName,
		_sessionRedisProviderMaxLifeTime,
		session.NewManager,
	)
	_sessionRedisProvider = wire.NewSet(
		redisProvider.NewRedisProvider,
		wire.Bind(new(session.Provider), new(*redisProvider.RedisProvider)),
	)
)

func permissionsHandler() (model.PermissionsHandler, error) {
	return model.NewPermissionsHandler("resource_action_rules.json")
}

var (
	entityHandler = wire.NewSet(
		_registerEntities,
		entity.NewEntityHandler,
	)
	entityUseCaseHandler = wire.NewSet(
		entityUsecase.NewEntityUseCase,
	)
)

func _registerEntities() []entity.RegisterOpt[any] {
	return []entity.RegisterOpt[any]{{}}
}

func _wrapListQueryFunc[T any](fn func(db *gorm.DB, opt T) *gorm.DB) eGorm.ListQueryFunc {
	return func(db *gorm.DB, listQueryOption any) *gorm.DB {
		return fn(db, listQueryOption.(T))
	}
}
