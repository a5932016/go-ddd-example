package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/a5932016/go-ddd-example/config"
)

func NewMySQL(env config.Environment) (*gorm.DB, error) {
	host := env.MySQL.Host
	user := env.MySQL.User
	port := env.MySQL.Port
	password := env.MySQL.Password
	dbName := env.MySQL.DBName
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", user, password, host, port, dbName)

	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN: dsn,
	}), &gorm.Config{
		TranslateError: true,
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold:             200 * time.Millisecond, // Slow SQL threshold
				LogLevel:                  logger.Error,
				IgnoreRecordNotFoundError: true,
				Colorful:                  true,
			},
		),
	})
	if err != nil {
		return nil, errors.Wrap(err, "connect db")
	}

	return db, nil
}

// NewRedis new redis
func NewRedis(env config.Environment) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     env.Redis.Host + ":" + env.Redis.Port,
		Password: env.Redis.Password,
	})
	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, errors.Wrap(err, "redis ping")
	}

	return client, nil
}
