package redis

import (
	"context"
	"time"

	"github.com/a5932016/go-ddd-example/repository"
	"github.com/a5932016/go-ddd-example/singleton/session"
	"github.com/a5932016/go-ddd-example/util/locking"
)

func NewRedisProvider(memRepo repository.MemRepository, maxLifeTime session.MaxLifeTime) *RedisProvider {
	return &RedisProvider{
		memRepo:     memRepo,
		maxLifeTime: maxLifeTime,
	}
}

type RedisProvider struct {
	memRepo     repository.MemRepository
	maxLifeTime session.MaxLifeTime
}

func (rp *RedisProvider) SessionInit(sid string) (session.Session, error) {
	defer locking.Lock(context.Background(), sid)()

	// Calculate the expiration time based on maxLifeTime
	expiration := time.Duration(rp.maxLifeTime) * time.Second

	// Set the session key with expiration
	if err := rp.memRepo.HSetEx(sid, "init", 1, expiration); err != nil {
		return nil, err
	}

	return &RedisSession{sid: sid, memRepo: rp.memRepo}, nil
}

func (rp *RedisProvider) SessionRead(sid string) (session.Session, error) {
	defer locking.Lock(context.Background(), sid)()

	// Refresh the session expiration time
	expiration := time.Duration(rp.maxLifeTime) * time.Second
	exists, err := rp.memRepo.Expire(sid, expiration)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, session.ErrSessionNotExisted
	}

	return &RedisSession{sid: sid, memRepo: rp.memRepo}, nil
}

func (rp *RedisProvider) SessionDestroy(sid string) error {
	defer locking.Lock(context.Background(), sid)()

	_, err := rp.memRepo.Del(sid)
	return err
}

func (rp *RedisProvider) SessionGC(maxLifeTime session.MaxLifeTime) {
	// Not implemented for Redis as it has its own expiration mechanism
}
