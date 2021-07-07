package config

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tigorlazuardi/ridit-go/app/config"
)

var printConfigCMD = &cobra.Command{
	Use:     "print",
	Aliases: []string{"list", "ls"},
	Short:   "prints current configuration to stdout",
	Run: func(cmd *cobra.Command, args []string) {
		format, _ := cmd.Flags().GetString("type")
		val, err := config.FormatConfig(viper.GetString("profile"), format)
		if err == config.ErrNotSupported {
			logrus.
				WithField("given_format", format).
				WithField("supported_formats", "yaml,toml,json").Fatal(err)
		} else if err != nil {
			logrus.Fatal(err)
		}
		fmt.Println(string(val))
	},
}

func init() {
	printConfigCMD.Flags().StringP("type", "t", "toml", "output format to print. supported types: [toml, yaml, json]")
	ConfigCMD.AddCommand(printConfigCMD)
}
