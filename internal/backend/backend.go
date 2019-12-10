package backend

import (
	"github.com/purpledb/purple"
	"github.com/purpledb/purple/internal/backend/disk"
	"github.com/purpledb/purple/internal/backend/memory"
	"github.com/purpledb/purple/internal/backend/redis"
	"github.com/purpledb/purple/services/cache"
	"github.com/purpledb/purple/services/counter"
	"github.com/purpledb/purple/services/flag"
	"github.com/purpledb/purple/services/kv"
	"github.com/purpledb/purple/services/set"
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

func NewBackend(cfg *purple.ServerConfig) (*Backend, error) {
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
		return nil, purple.ErrBackendNotRecognized
	}
}

func (b *Backend) Close() error {
	return b.Service.Close()
}
