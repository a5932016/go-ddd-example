package migration

import (
	"github.com/a5932016/go-ddd-example/repository"
	"github.com/a5932016/go-ddd-example/repository/casbin"
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

var migrations = []*gormigrate.Migration{
	firstMigration,
}

// New new migration
func New(
	dbRepo repository.DBRepository,
	perRepo *casbin.PERRepository,
) Migration {
	return Migration{
		dbRepo:  dbRepo,
		perRepo: perRepo,
	}
}

type Migration struct {
	dbRepo           repository.DBRepository
	perRepo          *casbin.PERRepository
	hierarchyPerRepo *casbin.PERRepository
}

func (m *Migration) Migrate() {
	m.dbRepo.Migrate(func(db *gorm.DB) error {
		gm := gormigrate.New(db, gormigrate.DefaultOptions, migrations)
		if err := gm.Migrate(); err != nil {
			return errors.Wrap(err, "gm.Migrate")
		}
		return nil
	})
}
