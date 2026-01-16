package redis

import (
	"fmt"

	"github.com/a5932016/go-ddd-example/repository"
)

type RedisSession struct {
	sid     string
	memRepo repository.MemRepository
}

func (rs *RedisSession) Set(key, value interface{}) error {
	_, err := rs.memRepo.HSet(rs.sid, key, value)
	return err
}

func (rs *RedisSession) Get(key interface{}) interface{} {
	v, _ := rs.memRepo.HGet(rs.sid, fmt.Sprintf("%v", key))
	return v
}

func (rs *RedisSession) Delete(key interface{}) error {
	_, err := rs.memRepo.HDel(rs.sid, fmt.Sprintf("%v", key))
	return err
}

func (rs *RedisSession) SessionID() string {
	return rs.sid
}
