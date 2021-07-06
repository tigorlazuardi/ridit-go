package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/tigorlazuardi/ridit-go/app/config"
)

var printConfigCMD = &cobra.Command{
	Use:     "print",
	Aliases: []string{"list", "ls"},
	Short:   "prints current configuration to stdout",
	Run: func(cmd *cobra.Command, args []string) {
		format, _ := cmd.Flags().GetString("type")
		val, err := config.FormatConfig(format)
		if err == config.ErrNotSupported {
			logrus.WithError(err).WithField("supported_formats", "yaml,toml,json").Fatal(err)
		} else if err != nil {
			logrus.Fatal(err)
		}
		logrus.Println(string(val))
	},
}

func init() {
	printConfigCMD.Flags().StringP("type", "t", "toml", "output format to print. supported types: [toml, yaml, json]")
	ConfigCMD.AddCommand(printConfigCMD)
}
