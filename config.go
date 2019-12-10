package purple

type ServerConfig struct {
	Port     int
	Debug    bool
	Backend  string
	RedisUrl string
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
