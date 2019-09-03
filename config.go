package strato

type (
	ClientConfig struct {
		Address string
	}

	ServerConfig struct {
		Port    int
		Debug   bool
		Backend string
	}
)

func (c *ClientConfig) Validate() error {
	if c.Address == "" {
		return ErrNoAddress
	}

	return nil
}

func (c *ServerConfig) Validate() error {
	if c.Port == 0 {
		return ErrNoPort
	}

	if c.Port < 1024 || c.Port > 49151 {
		return ErrPortOutOfRange
	}

	if c.Backend == "" {
		return ErrNoBackend
	}

	return nil
}
