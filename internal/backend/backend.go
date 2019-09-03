package backend

import (
	"github.com/lucperkins/strato"
	"github.com/lucperkins/strato/internal/backend/disk"
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
	_ Backend = (*strato.Memory)(nil)
)

func NewBackend(cfg *strato.ServerConfig) (*Holder, error) {
	switch cfg.Backend {
	case "disk":
		backend, err := disk.NewDisk()
		if err != nil {
			return nil, err
		}
		return &Holder{
			backend,
		}, nil
	case "memory":
		return &Holder{
			strato.NewMemoryBackend(),
		}, nil
	default:
		return nil, strato.ErrBackendNotRecognized
	}
}

func (b *Holder) Close() error {
	return b.Backend.Close()
}
