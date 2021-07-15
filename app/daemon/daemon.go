package daemon

import (
	"github.com/spf13/cobra"
	"github.com/tigorlazuardi/ridit/pkg"
)

var daemonCMD = &cobra.Command{
	Use:     "daemon",
	Aliases: []string{"http", "server"},
	Short:   "runs the program as daemon",
	Run: func(cmd *cobra.Command, args []string) {
		entry := pkg.EntryFromContext(cmd.Context())
		entry.Fatal("not implemented")
	},
}
