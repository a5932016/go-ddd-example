package redis

import (
	"time"

	"github.com/go-redsync/redsync/v4"
)

var defaultMutexOpts = []redsync.Option{
	redsync.WithRetryDelay(3 * time.Second),
	redsync.WithExpiry(3 * time.Minute),
}

// NewMutex new mutex
func (m *MemRepository) NewMutex(name string, opts ...redsync.Option) *redsync.Mutex {
	options := append(defaultMutexOpts, opts...)
	return m.redsync.NewMutex(name, options...)
}
