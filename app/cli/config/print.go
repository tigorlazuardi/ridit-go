package config

import "github.com/spf13/cobra"

var printConfigCMD = &cobra.Command{
	Use:     "print",
	Aliases: []string{"list", "ls"},
	Short:   "prints current configuration to stdout",
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	ConfigCMD.AddCommand(printConfigCMD)
}
