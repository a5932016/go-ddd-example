package main

import (
	"os"

	"github.com/a5932016/go-ddd-example/config"
	"github.com/a5932016/go-ddd-example/db"
	"github.com/a5932016/go-ddd-example/util/log"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/afero"
	"github.com/urfave/cli"
	"gorm.io/gorm"
)

var (
	app           *cli.App
	rollbackSteps int
	envPath       string

	redisC    *redis.Client
	postgresC *gorm.DB

	// Version control.
	Version      = "No Version Provided"
	BuildDate    = ""
	GitCommitSha = ""
)

func init() {
	app = cli.NewApp()
	app.Name = "github.com/a5932016/go-ddd-example"
	app.Version = "1"

	app.Action = func(c *cli.Context) error {
		if err := config.InitEnvironment(envPath); err != nil {
			return errors.Wrap(err, "init env")
		}
		// Redis
		redisC, err := db.NewRedis(config.Env)
		if err != nil {
			return errors.Wrap(err, "db.NewRedis")
		}
		// MySQL
		mySqlC, err := db.NewMySQL(config.Env)
		if err != nil {
			return errors.Wrap(err, "db.NewMySQL")
		}

		// File System
		appFs := afero.NewOsFs()

		// Migration
		migration, err := InitMigration(mySqlC)
		if err != nil {
			return errors.Wrap(err, "InitMigration")
		}

		migration.Migrate()

		// init router
		router, err := InitRouter(appFs, mySqlC, redisC)
		if err != nil {
			return errors.Wrap(err, "InitRouter")
		}

		if err := router.RunServer(); err != nil {
			return errors.Wrap(err, "router.RunServer")
		}
		return nil
	}
}

func main() {
	if err := app.Run(os.Args); err != nil {
		log.WithError(err).Fatal("Service Run")
	}
}
