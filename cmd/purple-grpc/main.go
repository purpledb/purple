package main

import (
	"github.com/purpledb/purple"
	"github.com/purpledb/purple/cmd"
	"github.com/purpledb/purple/internal/server/grpc"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func command() *cobra.Command {
	var cfg purple.ServerConfig

	v := cmd.NewConfig("purple_grpc")

	command := &cobra.Command{
		Use: "purple-grpc",
		PreRun: func(_ *cobra.Command, _ []string) {
			cmd.ExitOnError(v.Unmarshal(&cfg))
		},
		Run: func(_ *cobra.Command, _ []string) {
			srv, err := grpc.NewGrpcServer(&cfg)
			cmd.ExitOnError(err)
			cmd.ExitOnError(srv.Start())
		},
	}

	flags := pflag.NewFlagSet("purple-grpc", pflag.ExitOnError)
	flags.IntP("port", "p", 8081, "Purple server port")
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
