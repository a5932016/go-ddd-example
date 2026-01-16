package mysql

import (
	"time"

	"github.com/a5932016/go-ddd-example/repository"
	"github.com/a5932016/go-ddd-example/singleton/entity"
	"github.com/a5932016/go-ddd-example/util/log"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	defaultLockTimeout     = 1 * time.Minute
	defaultTxEffectiveTime = 3 * time.Minute
)

// NewDBRepository func implements the storage interface for app
func NewDBRepository(db *gorm.DB, entityCtrl *entity.EntityHandler) *DBRepository {
	repo := &DBRepository{
		db:         db,
		entityCtrl: entityCtrl,
	}
	repo.SetDefaultLogMode()
	return repo
}

// DBRepository is interface structure
type DBRepository struct {
	db         *gorm.DB
	entityCtrl *entity.EntityHandler
}

// DB
func (s *DBRepository) DB() *gorm.DB {
	return s.db
}

func (s *DBRepository) SetDefaultLogMode() {
	s.db.Logger = s.db.Logger.LogMode(logger.Error)
}

func (s *DBRepository) SetLogMode(mode logger.LogLevel) {
	s.db.Logger = s.db.Logger.LogMode(mode)
}

func (m *DBRepository) Migrate(fn func(*gorm.DB) error) {
	m.SetLogMode(logger.Info)
	log.Info("Start Migration")

	if err := fn(m.db); err != nil {
		log.WithError(err).Error("Database Migration Failed")
	}

	log.Info("End Migration")
	m.SetDefaultLogMode()
}

func (s *DBRepository) EntityCtrl() *entity.EntityHandler {
	return s.entityCtrl
}

// Begin begin a transaction
func (s *DBRepository) Begin() repository.DBRepository {
	return &DBRepository{
		db: s.db.Begin(),
	}
}

// Commit commit a transaction
func (s *DBRepository) Commit() error {
	return s.db.Commit().Error
}

// Rollback rollback a transaction
func (s *DBRepository) Rollback() error {
	return s.db.Rollback().Error
}

// Debug Debug log
func (s *DBRepository) Debug() {
	s.db = s.db.Debug()
	return
}
