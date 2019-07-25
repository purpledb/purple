package server

type Config struct {
	Port int
}

func (c *Config) validate() error {
	if c.Port == 0 {
		return ErrNoPort
	}

	if c.Port < 1024 || c.Port > 49151 {
		return ErrPortOutOfRange
	}

	return nil
}
