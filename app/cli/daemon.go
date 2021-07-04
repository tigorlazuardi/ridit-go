package cli

import (
	"github.com/spf13/cobra"
)

var daemonCMD = &cobra.Command{
	Use:     "daemon",
	Aliases: []string{"http", "serve"},
	Short:   "runs ridit as daemon http service",
	Long:    "runs ridit as daemon http service",
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(daemonCMD)
}
