package config

import "github.com/lucperkins/strato"

type (
	ClientConfig struct {
		Address string
	}

	ServerConfig struct {
		Port     int
		Debug    bool
		Backend  string
		RedisUrl string
	}
)

func (c *ClientConfig) Validate() error {
	if c.Address == "" {
		return strato.ErrNoAddress
	}

	return nil
}

func (c *ServerConfig) Validate() error {
	if c.Port == 0 {
		return strato.ErrNoPort
	}

	if c.Port < 1024 || c.Port > 49151 {
		return strato.ErrPortOutOfRange
	}

	if c.Backend == "" {
		return strato.ErrNoBackend
	}

	return nil
}
