package backend

import (
	"github.com/lucperkins/strato"
	"github.com/lucperkins/strato/internal/backend/disk"
	"github.com/lucperkins/strato/internal/backend/memory"
	"github.com/lucperkins/strato/internal/backend/redis"
	"github.com/lucperkins/strato/internal/services/cache"
	"github.com/lucperkins/strato/internal/services/counter"
	"github.com/lucperkins/strato/internal/services/flag"
	"github.com/lucperkins/strato/internal/services/kv"
	"github.com/lucperkins/strato/internal/services/set"
)

type (
	Service interface {
		cache.Cache
		counter.Counter
		flag.Flag
		kv.KV
		set.Set

		Close() error
		Flush() error
		Name() string
	}

	// Backend wraps a Service and thereby provides specific instantiations access to the Close() and Flush() methods
	Backend struct {
		Service
	}
)

var (
	_ Service = (*disk.Disk)(nil)
	_ Service = (*memory.Memory)(nil)
	_ Service = (*redis.Redis)(nil)
)

func NewBackend(cfg *strato.ServerConfig) (*Backend, error) {
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
	return b.Service.Close()
}
