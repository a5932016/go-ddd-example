package repository

import (
	"time"

	"github.com/go-redsync/redsync/v4"
	"github.com/ulule/limiter/v3"
)

// MemRepository interface
type MemRepository interface {
	redisI
}

type redisI interface {
	NewMutex(name string, options ...redsync.Option) *redsync.Mutex
	GetAPILimiter() limiter.Store

	Set(key string, value interface{}, expiration time.Duration) (string, error)
	Del(keys ...string) (int64, error)

	HSet(key string, values ...interface{}) (int64, error)
	HSetEx(key, field string, value interface{}, expiration time.Duration) error
	HGet(key string, field string) (string, error)
	HDel(key string, fields ...string) (int64, error)

	Expire(key string, expiration time.Duration) (bool, error)
	Exists(keys ...string) (int64, error)
}
