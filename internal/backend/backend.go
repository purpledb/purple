package backend

import (
	"github.com/lucperkins/strato"
	"github.com/lucperkins/strato/internal/backend/disk"
	"github.com/lucperkins/strato/internal/backend/memory"
)

type (
	Interface interface {
		strato.Cache
		strato.Counter
		strato.KV
		strato.Set

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
	default:
		return nil, strato.ErrBackendNotRecognized
	}
}

func (b *Backend) Close() error {
	return b.Interface.Close()
}
