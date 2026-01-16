package eGorm

import (
	"fmt"
)

func NewErrorEntityGORM[T Entity](entity T) ErrorEntityGORM[T] {
	return ErrorEntityGORM[T]{
		entity: entity,
	}
}

type ErrorEntityGORM[T Entity] struct {
	entity Entity
}

func (e ErrorEntityGORM[T]) Get(dist any, id uint) error {
	return errFn(e.entity)
}

func (e ErrorEntityGORM[T]) List(dist any, opt any) error {
	return errFn(e.entity)
}

func (e ErrorEntityGORM[T]) Count(dist *int64, opt any) error {
	return errFn(e.entity)
}

func (e ErrorEntityGORM[T]) Create(dist any) error {
	return errFn(e.entity)
}

func (e ErrorEntityGORM[T]) Update(dist any, id uint) error {
	return errFn(e.entity)
}

func (e ErrorEntityGORM[T]) Delete(id uint) error {
	return errFn(e.entity)
}

var errFn = func(entity Entity) error {
	return fmt.Errorf("Entity %s isn't registered", entity.ModelName())
}
