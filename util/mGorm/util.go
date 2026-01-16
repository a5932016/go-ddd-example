package mGorm

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var LockingStrengthShare = clause.Locking{
	Strength: clause.LockingStrengthShare,
	Table:    clause.Table{Name: clause.CurrentTable},
}

func GetIDs[T any](models []T, idSelector func(T) uint) []uint {
	ids := make([]uint, len(models))
	for _, model := range models {
		ids = append(ids, idSelector(model))
	}
	return ids
}

func FindIdDiff[T any](oldItems []T, newItems []T, idSelector func(T) uint) ([]uint, []uint) {
	oldSet := make(map[uint]bool, len(oldItems))
	for _, item := range oldItems {
		oldSet[idSelector(item)] = true
	}
	newSet := make(map[uint]bool, len(newItems))
	for _, item := range newItems {
		newSet[idSelector(item)] = true
	}

	// Pre-allocate result sizes for potential improvement
	oldOnly := make([]uint, 0, len(oldItems))
	newOnly := make([]uint, 0, len(newItems))

	// Find Differences
	for _, item := range oldItems {
		if !newSet[idSelector(item)] {
			oldOnly = append(oldOnly, idSelector(item))
		}
	}
	for _, item := range newItems {
		if !oldSet[idSelector(item)] {
			newOnly = append(newOnly, idSelector(item))
		}
	}
	return oldOnly, newOnly
}

// Generic batch create function
func BatchCreate[T any](db *gorm.DB, items []T, batchSize int) error {
	if batchSize <= 0 {
		return errors.New("batchSize must be greater than 0")
	}

	for i := 0; i < len(items); i += batchSize {
		end := i + batchSize
		if end > len(items) {
			end = len(items)
		}
		if err := db.CreateInBatches(items[i:end], batchSize).Error; err != nil {
			return errors.Wrapf(err, "failed to create batch %d-%d", i, end)
		}
	}
	return nil
}

func CheckSliceLengths[T any](slice1 []T, slice2 []T) bool {
	return len(slice1) != len(slice2)
}
