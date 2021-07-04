package config

import (
	"github.com/spf13/cobra"
)

var ConfigCMD = &cobra.Command{
	Use:     "config",
	Short:   "set or get configurations",
	Aliases: []string{"cfg"},
	Example: "ridit config set aspect_ratio enable",
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}
