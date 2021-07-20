package cli

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/tigorlazuardi/ridit/app/daemon/server"
	"github.com/tigorlazuardi/ridit/app/daemon/server/router"
	"github.com/tigorlazuardi/ridit/pkg"
)

var daemonCMD = &cobra.Command{
	Use:     "daemon",
	Aliases: []string{"http", "serve"},
	Short:   "runs ridit as daemon http service",
	Long:    "runs ridit as daemon http service",
	Run: func(cmd *cobra.Command, _ []string) {
		entry := pkg.EntryFromContext(cmd.Context())
		sig := pkg.RegisterInterrupt()

		router := router.New()

		// TODO: move this 10101 to config port
		close := server.Start(router, ":10101")
		log.Println("http server running on :10101")
		<-sig
		err := close()
		if err != nil {
			entry.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(daemonCMD)
}
