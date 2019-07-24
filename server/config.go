package server

type Config struct {
	Port int
}

func (c *Config) validate() error {
	if c.Port == 0 {
		return ErrNoPort()
	}

	return nil
}