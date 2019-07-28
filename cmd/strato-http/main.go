package main

import (
	"github.com/lucperkins/strato"
	"github.com/lucperkins/strato/cmd"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func command() *cobra.Command {
	var config strato.HttpConfig

	v := viper.New()
	v.AutomaticEnv()
	v.SetEnvPrefix("strato_http")

	command := &cobra.Command{
		Use: "strato-http",
		PreRun: func(_ *cobra.Command, _ []string) {
			cmd.ExitOnError(v.Unmarshal(&config))
		},
		Run: func(_ *cobra.Command, _ []string) {
			srv := strato.NewHttpServer(&config)
			cmd.ExitOnError(srv.Start())
		},
	}

	flags := pflag.NewFlagSet("strato-http", pflag.ExitOnError)
	flags.IntP("port", "p", 8081, "Strato HTTP server port")
	flags.Bool("debug", false, "Debug mode")

	cmd.BindFlagsToCmd(command, flags, v)

	return command
}
func main() {
	cmd.ExitOnError(command().Execute())
}
