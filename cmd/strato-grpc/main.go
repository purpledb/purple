package main

import (
	"github.com/lucperkins/strato/cmd"
	"github.com/lucperkins/strato/internal/config"
	"github.com/lucperkins/strato/internal/server/grpc"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func command() *cobra.Command {
	var cfg config.ServerConfig

	v := viper.New()
	v.AutomaticEnv()
	v.SetEnvPrefix("strato_grpc")

	command := &cobra.Command{
		Use: "strato-grpc",
		PreRun: func(_ *cobra.Command, _ []string) {
			cmd.ExitOnError(v.Unmarshal(&cfg))
		},
		Run: func(_ *cobra.Command, _ []string) {
			srv, err := grpc.NewGrpcServer(&cfg)
			cmd.ExitOnError(err)
			cmd.ExitOnError(srv.Start())
		},
	}

	flags := pflag.NewFlagSet("strato-grpc", pflag.ExitOnError)
	flags.IntP("port", "p", 8080, "Strato server port")
	flags.Bool("debug", false, "Debug mode")
	flags.String("backend", "disk", `Data backend (options are "disk" and "memory")`)
	flags.String("redis-url", "localhost:6379", "Redis connection URL (if redis backend is used)")

	v.RegisterAlias("redisurl", "redis-url")

	cmd.BindFlagsToCmd(command, flags, v)

	return command
}

func main() {
	cmd.ExitOnError(command().Execute())
}
