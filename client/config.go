package client

type Config struct {
	Address string
}

func (c *Config) validate() error {
	if c.Address == "" {
		return ErrNoAddress
	}

	return nil
}
