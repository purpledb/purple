package main

import (
	"github.com/lucperkins/strato"
	"github.com/lucperkins/strato/cmd"
	"github.com/lucperkins/strato/internal/server/http"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func command() *cobra.Command {
	var cfg strato.ServerConfig

	v := viper.New()
	v.AutomaticEnv()
	v.SetEnvPrefix("strato_http")

	command := &cobra.Command{
		Use: "strato-http",
		PreRun: func(_ *cobra.Command, _ []string) {
			cmd.ExitOnError(v.Unmarshal(&cfg))
		},
		Run: func(_ *cobra.Command, _ []string) {
			srv, err := http.NewHttpServer(&cfg)
			cmd.ExitOnError(err)
			cmd.ExitOnError(srv.Start())
		},
	}

	flags := pflag.NewFlagSet("strato-http", pflag.ExitOnError)
	flags.IntP("port", "p", 8081, "Strato HTTP server port")
	flags.Bool("debug", false, "Debug mode")
	flags.String("backend", "disk", `Data backend (options are "disk" and "memory")`)

	cmd.BindFlagsToCmd(command, flags, v)

	return command
}
func main() {
	cmd.ExitOnError(command().Execute())
}
