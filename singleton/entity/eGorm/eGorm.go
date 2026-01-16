package eGorm

import (
	"gorm.io/gorm"
)

type Entity interface {
	ModelName() string
}

type ListQueryFunc func(*gorm.DB, any) *gorm.DB

func NewEntityGORM[T Entity](db *gorm.DB, entity T, listQueryFn ListQueryFunc) EntityGORM[T] {
	return EntityGORM[T]{
		db:          db,
		entity:      entity,
		listQueryFn: listQueryFn,
	}
}

type EntityGORM[T Entity] struct {
	db          *gorm.DB
	entity      Entity
	listQueryFn ListQueryFunc
}

func (e EntityGORM[T]) Get(dist any, id uint) error {
	return e.db.Where("id = ?", id).First(dist).Error
}

func (e EntityGORM[T]) List(dist any, opt any) error {
	return e.getDB(opt).Find(dist).Error
}

func (e EntityGORM[T]) Count(dist *int64, opt any) error {
	return e.getDB(opt).Model(e.entity).Select("*").Limit(-1).Offset(-1).Count(dist).Error
}

func (e EntityGORM[T]) getDB(opt any) *gorm.DB {
	if e.listQueryFn == nil {
		return e.db
	}
	return e.listQueryFn(e.db, opt)
}

func (e EntityGORM[T]) Create(dist any) error {
	return e.db.Create(dist).Error
}

func (e EntityGORM[T]) Update(dist any, id uint) error {
	return e.db.Where("id = ?", id).Updates(dist).Error
}

func (e EntityGORM[T]) Delete(id uint) error {
	return e.db.Where("id = ?", id).Delete(&e.entity).Error
}
