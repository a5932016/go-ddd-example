package redis

import (
	"context"
	"sync"
	"time"

	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"github.com/ulule/limiter/v3"
	sredis "github.com/ulule/limiter/v3/drivers/store/redis"
)

// NewMemRepository func implements the storage interface for app
func NewMemRepository(redisClient *redis.Client) (*MemRepository, error) {
	limiterStore, err := sredis.NewStoreWithOptions(redisClient, limiter.StoreOptions{
		Prefix:   "fe_api_rate_limit",
		MaxRetry: 3,
	})
	if err != nil {
		return nil, errors.Wrap(err, "sredis.NewStoreWithOptions")
	}
	return &MemRepository{
		client:     redisClient,
		redsync:    redsync.New(goredis.NewPool(redisClient)),
		apiLimiter: &limiterStore,
		ctx:        context.Background(),
	}, nil
}

// MemRepository is interface structure
type MemRepository struct {
	client      *redis.Client
	redsync     *redsync.Redsync
	apiLimiter  *limiter.Store
	scriptSHA1s sync.Map
	ctx         context.Context
}

func (m *MemRepository) GetAPILimiter() limiter.Store {
	return *m.apiLimiter
}

func (m *MemRepository) Set(key string, value interface{}, expiration time.Duration) (string, error) {
	return m.client.Set(m.ctx, key, value, expiration).Result()
}

func (m *MemRepository) Del(keys ...string) (int64, error) {
	return m.client.Del(m.ctx, keys...).Result()
}

func (m *MemRepository) HSet(key string, values ...interface{}) (int64, error) {
	return m.client.HSet(m.ctx, key, values...).Result()
}

func (m *MemRepository) HSetEx(key, field string, value interface{}, expiration time.Duration) error {
	script := redis.NewScript(`
        local key = KEYS[1]
        local field = ARGV[1]
        local value = ARGV[2]
        local ttl = tonumber(ARGV[3])

        local result = redis.call('HSET', key, field, value)
        if result == 1 then
          redis.call('EXPIRE', key, ttl)
          return 'OK'
        else
          return 'Error: Could not set field'
        end
    `)
	_, err := script.Run(m.ctx, m.client, []string{key}, field, value, formatSec(expiration)).Result()
	return err
}

func (m *MemRepository) HGet(key string, field string) (string, error) {
	return m.client.HGet(m.ctx, key, field).Result()
}

func (m *MemRepository) HDel(key string, fields ...string) (int64, error) {
	return m.client.HDel(m.ctx, key, fields...).Result()
}

func (m *MemRepository) Expire(key string, expiration time.Duration) (bool, error) {
	return m.client.Expire(m.ctx, key, expiration).Result()
}

func (m *MemRepository) Exists(keys ...string) (int64, error) {
	return m.client.Exists(m.ctx, keys...).Result()
}

func formatSec(dur time.Duration) int64 {
	if dur > 0 && dur < time.Second {
		return 1
	}
	return int64(dur / time.Second)
}
