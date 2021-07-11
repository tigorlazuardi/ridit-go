package cli

import (
	"github.com/spf13/cobra"
	"github.com/tigorlazuardi/ridit/pkg"
)

var daemonCMD = &cobra.Command{
	Use:     "daemon",
	Aliases: []string{"http", "serve"},
	Short:   "runs ridit as daemon http service",
	Long:    "runs ridit as daemon http service",
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
		entry := pkg.EntryFromContext(cmd.Context())
		entry.Fatal("not implemented")
	},
}

func init() {
	rootCmd.AddCommand(daemonCMD)
}
