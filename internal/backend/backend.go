package backend

import (
	"github.com/lucperkins/strato"
	"github.com/lucperkins/strato/internal/backend/disk"
	"github.com/lucperkins/strato/internal/backend/memory"
)

type (
	Backend interface {
		strato.Cache
		strato.Counter
		strato.KV
		strato.Set

		Close() error
		Flush() error
	}

	Holder struct {
		Backend
	}
)

var (
	_ Backend = (*disk.Disk)(nil)
	_ Backend = (*memory.Memory)(nil)
)

func NewBackend(cfg *strato.ServerConfig) (*Holder, error) {
	switch cfg.Backend {
	case "disk":
		backend, err := disk.NewDiskBackend()
		if err != nil {
			return nil, err
		}
		return &Holder{
			backend,
		}, nil
	case "memory":
		backend := memory.NewMemoryBackend()

		return &Holder{
			backend,
		}, nil
	default:
		return nil, strato.ErrBackendNotRecognized
	}
}

func (b *Holder) Close() error {
	return b.Backend.Close()
}
