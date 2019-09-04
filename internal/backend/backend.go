package backend

import (
	"github.com/lucperkins/strato"
	"github.com/lucperkins/strato/internal/backend/disk"
	"github.com/lucperkins/strato/internal/backend/memory"
	"github.com/lucperkins/strato/internal/backend/redis"
	"github.com/lucperkins/strato/internal/config"
	"github.com/lucperkins/strato/internal/services/cache"
	"github.com/lucperkins/strato/internal/services/counter"
	"github.com/lucperkins/strato/internal/services/kv"
	"github.com/lucperkins/strato/internal/services/set"
)

type (
	Interface interface {
		cache.Cache
		counter.Counter
		kv.KV
		set.Set

		Close() error
		Flush() error
	}

	Backend struct {
		Interface
	}
)

var (
	_ Interface = (*disk.Disk)(nil)
	_ Interface = (*memory.Memory)(nil)
	_ Interface = (*redis.Redis)(nil)
)

func NewBackend(cfg *config.ServerConfig) (*Backend, error) {
	switch cfg.Backend {
	case "disk":
		backend, err := disk.NewDiskBackend()
		if err != nil {
			return nil, err
		}
		return &Backend{
			backend,
		}, nil
	case "memory":
		backend := memory.NewMemoryBackend()

		return &Backend{
			backend,
		}, nil
	case "redis":
		backend, err := redis.NewRedisBackend(cfg.RedisUrl)
		if err != nil {
			return nil, err
		}
		return &Backend{
			backend,
		}, nil
	default:
		return nil, strato.ErrBackendNotRecognized
	}
}

func (b *Backend) Close() error {
	return b.Interface.Close()
}
