package strato

type (
	ClientConfig struct {
		Address string
	}

	ServerConfig struct {
		Port int
	}
)

func (c *ClientConfig) validate() error {
	if c.Address == "" {
		return ErrNoAddress
	}

	return nil
}

func (c *ServerConfig) validate() error {
	if c.Port == 0 {
		return ErrNoPort
	}

	if c.Port < 1024 || c.Port > 49151 {
		return ErrPortOutOfRange
	}

	return nil
}
