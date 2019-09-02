package strato

type (
	Backend interface {
		Cache
		Counter
		KV
		Set

		Close() error
		Flush() error
	}

	BackendObj struct {
		Backend
	}
)

var (
	_ Backend = (*Disk)(nil)
	_ Backend = (*Memory)(nil)
)

func NewBackend(cfg *ServerConfig) (*BackendObj, error) {
	switch cfg.Backend {
	case "disk":
		backend, err := NewDisk(dbDir)
		if err != nil {
			return nil, err
		}
		return &BackendObj{
			backend,
		}, nil
	case "memory":
		return &BackendObj{
			NewMemoryBackend(),
		}, nil
	default:
		return nil, ErrBackendNotRecognized
	}
}

func (b *BackendObj) Close() error {
	return b.Backend.Close()
}
