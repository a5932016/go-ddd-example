package entity

import (
	"gorm.io/gorm"

	"github.com/a5932016/go-ddd-example/singleton/entity/eGorm"
)

type RegisterOpt[T any] struct {
	Entity        eGorm.Entity
	ListQueryFunc eGorm.ListQueryFunc
}

func NewEntityHandler(db *gorm.DB, opts []RegisterOpt[any]) *EntityHandler {
	entityGORMs := make(map[string]eGorm.EntityGORM[eGorm.Entity])
	for _, opt := range opts {
		entityGORMs[opt.Entity.ModelName()] = eGorm.NewEntityGORM(db, opt.Entity, opt.ListQueryFunc)
	}
	return &EntityHandler{
		opts:        opts,
		entityGORMs: entityGORMs,
	}
}

type EntityHandler struct {
	opts        []RegisterOpt[any]
	entityGORMs map[string]eGorm.EntityGORM[eGorm.Entity]
}

func (h *EntityHandler) Entity(entity eGorm.Entity) EGORM {
	entityGORM, ok := h.entityGORMs[entity.ModelName()]
	if ok {
		return entityGORM
	}

	return eGorm.NewErrorEntityGORM(entity)
}

func (h *EntityHandler) Begin(db *gorm.DB) *EntityHandler {
	return NewEntityHandler(db, h.opts)
}
