package main

import (
	"github.com/google/wire"
	"github.com/spf13/afero"
	"gorm.io/gorm"

	"github.com/a5932016/go-ddd-example/repository"
	"github.com/a5932016/go-ddd-example/repository/casbin"
	"github.com/a5932016/go-ddd-example/repository/fs"
	"github.com/a5932016/go-ddd-example/repository/mysql"
	"github.com/a5932016/go-ddd-example/repository/redis"
)

var repositoryProvider = wire.NewSet(
	memRepoProvider,
	dbRepoProvider,
	perRepoProvider,
	fsRepoProvider,
)

var (
	memRepoProvider = wire.NewSet(
		redis.NewMemRepository,
		wire.Bind(new(repository.MemRepository), new(*redis.MemRepository)),
	)
	dbRepoProvider = wire.NewSet(
		mysql.NewDBRepository,
		wire.Bind(new(repository.DBRepository), new(*mysql.DBRepository)),
	)
)

func perRepoProvider(db *gorm.DB) (*casbin.PERRepository, error) {
	return casbin.NewPERRepository(db, "casbin.conf", "casbin_rules")
}

func fsRepoProvider(fsLib afero.Fs) fs.FSRepository {
	return fs.NewFSRepository(fsLib)
}
