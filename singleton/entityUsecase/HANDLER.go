package entityUsecase

import (
	"context"
	"fmt"

	"github.com/a5932016/go-ddd-example/customerror"
	"github.com/a5932016/go-ddd-example/repository"
	"github.com/a5932016/go-ddd-example/singleton/entity/eGorm"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

func NewEntityUseCase(dbRepo repository.DBRepository) EntityUseCase {
	return EntityUseCase{
		dbRepo: dbRepo,
	}
}

type EntityUseCase struct {
	dbRepo repository.DBRepository
}

func (h EntityUseCase) List(c context.Context, entity eGorm.Entity, opt any, dist any) (total int64, err error) {
	tx := h.dbRepo.Begin()
	defer tx.Rollback()
	entityGORM := tx.EntityCtrl().Entity(entity)
	if err = entityGORM.List(dist, opt); err != nil {
		return 0, errors.Wrap(err, fmt.Sprintf("Entity(%s).List", entity.ModelName()))
	}
	if err = entityGORM.Count(&total, opt); err != nil {
		return 0, errors.Wrap(err, fmt.Sprintf("Entity(%s).Count", entity.ModelName()))
	}
	return
}

func (h EntityUseCase) Create(c context.Context, entity eGorm.Entity, dist any) error {
	if err := h.dbRepo.EntityCtrl().Entity(entity).Create(dist); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return customerror.DuplicateName
		}
		return errors.Wrap(err, fmt.Sprintf("Entity(%s).Create", entity.ModelName()))
	}
	return nil
}

func (h EntityUseCase) Update(c context.Context, entity eGorm.Entity, id uint, dist any) error {
	tx := h.dbRepo.Begin()
	defer tx.Rollback()
	entityGORM := tx.EntityCtrl().Entity(entity)
	if err := entityGORM.Update(dist, id); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return customerror.DuplicateName
		}
		return errors.Wrap(err, fmt.Sprintf("Entity(%s).Update", entity.ModelName()))
	}
	if err := entityGORM.Get(dist, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return customerror.RecordNotFound
		}
		return errors.Wrap(err, fmt.Sprintf("Entity(%s).Get", entity.ModelName()))
	}
	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "commit")
	}
	return nil
}

func (h EntityUseCase) Delete(c context.Context, entity eGorm.Entity, id uint) error {
	if err := h.dbRepo.EntityCtrl().Entity(entity).Delete(id); err != nil {
		return errors.Wrap(err, fmt.Sprintf("Entity(%s).Delete", entity.ModelName()))
	}
	return nil
}
