package casbin

import (
	"context"
	"sync"

	"github.com/a5932016/go-ddd-example/util/log"
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

func NewPERRepository(db *gorm.DB, configPath, tableName string) (*PERRepository, error) {
	enforcer, err := newEnforcer(db, configPath, tableName)
	if err != nil {
		return nil, err
	}

	return &PERRepository{
		enforcer:   enforcer,
		lock:       new(sync.RWMutex),
		configPath: configPath,
		tableName:  tableName,
	}, nil
}

func newEnforcer(db *gorm.DB, configPath, tableName string) (*casbin.CachedEnforcer, error) {
	adapter, err := gormadapter.NewAdapterByDBUseTableName(db, "", tableName)
	if err != nil {
		return nil, errors.Wrap(err, "gormadapter.NewAdapterByDBUseTableName")
	}

	enforcer, err := casbin.NewCachedEnforcer(configPath, adapter) // will auto migrate
	if err != nil {
		return nil, errors.Wrap(err, "casbin.NewEnforcer")
	}

	return enforcer, nil
}

// PERRepository PERRepository
type PERRepository struct {
	enforcer   *casbin.CachedEnforcer
	withTx     bool
	lock       *sync.RWMutex
	configPath string
	tableName  string
}

func (r *PERRepository) BeginWithTx(db *gorm.DB) (txPerRepo *PERRepository, closeTx func(c context.Context), err error) {
	enforcer, err := newEnforcer(db, r.configPath, r.tableName)
	if err != nil {
		return nil, nil, err
	}

	closeTx = func(c context.Context) {
		r.lock.Unlock()
		if err := r.LoadPolicy(); err != nil {
			log.FromContext(c).WithError(err).Error("CloseTx")
		}
	}

	r.lock.Lock()

	return &PERRepository{
		enforcer:   enforcer,
		withTx:     true,
		lock:       r.lock,
		configPath: r.configPath,
		tableName:  r.tableName,
	}, closeTx, nil
}

func (r *PERRepository) Enforce(rvals ...interface{}) (bool, error) {
	if !r.withTx {
		r.lock.RLock()
		defer r.lock.RUnlock()
	}

	return r.enforcer.Enforce(rvals...)
}

func (r *PERRepository) LoadPolicy() error {
	if !r.withTx {
		r.lock.RLock()
		defer r.lock.RUnlock()
	}

	return r.enforcer.LoadPolicy()
}

func (r *PERRepository) GetPolicies(prefixedDivisionNameId string) ([][]string, error) {
	if !r.withTx {
		r.lock.RLock()
		defer r.lock.RUnlock()
	}

	return r.enforcer.GetPermissionsForUser(prefixedDivisionNameId)
}

func (r *PERRepository) AddPoliciesEx(rules [][]string) error {
	if !r.withTx {
		r.lock.Lock()
		defer r.lock.Unlock()
	}

	_, err := r.enforcer.AddPoliciesEx(rules)

	return err
}

func (r *PERRepository) RemovePolicies(rules [][]string) error {
	if !r.withTx {
		r.lock.Lock()
		defer r.lock.Unlock()
	}

	_, err := r.enforcer.RemovePolicies(rules)

	return err
}
