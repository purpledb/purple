package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"os"
)

func BindFlagsToCmd(cmd *cobra.Command, flags *pflag.FlagSet, v *viper.Viper) {
	cmd.Flags().AddFlagSet(flags)
	ExitOnError(v.BindPFlags(flags))
}

func ExitOnError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
