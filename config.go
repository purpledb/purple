package strato

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type (
	ClientConfig struct {
		Address string
	}

	GrpcConfig struct {
		Port  int
		Debug bool
	}

	HttpConfig struct {
		Port  int
		Debug bool
	}
)

func (c *ClientConfig) validate() error {
	if c.Address == "" {
		return ErrNoAddress
	}

	return nil
}

func (c *GrpcConfig) validate() error {
	if c.Port == 0 {
		return ErrNoPort
	}

	if c.Port < 1024 || c.Port > 49151 {
		return ErrPortOutOfRange
	}

	return nil
}

func getGrpcServerConfig(args []string) *GrpcConfig {
	var config GrpcConfig

	v := viper.New()
	v.AutomaticEnv()
	v.SetEnvPrefix("strato_grpc")

	flags := pflag.NewFlagSet("strato-grpc-server", pflag.ExitOnError)
	flags.IntP("port", "p", 8080, "Strato gRPC server port")
	flags.Bool("debug", false, "Debug mode")
	exitOnErr(flags.Parse(args))
	exitOnErr(v.BindPFlags(flags))
	exitOnErr(v.Unmarshal(&config))

	return &config
}

func getHttpServerConfig(args []string) *HttpConfig {
	var config HttpConfig

	v := viper.New()
	v.AutomaticEnv()
	v.SetEnvPrefix("strato_http")

	flags := pflag.NewFlagSet("strato-http-server", pflag.ExitOnError)
	flags.IntP("port", "p", 8081, "Strato HTTP server port")
	flags.Bool("debug", false, "Debug mode")
	exitOnErr(flags.Parse(args))
	exitOnErr(v.BindPFlags(flags))
	exitOnErr(v.Unmarshal(&config))

	return &config
}
