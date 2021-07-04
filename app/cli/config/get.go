package config

import (
	"github.com/spf13/cobra"
	aspectratio "github.com/tigorlazuardi/ridit-go/app/cli/config/aspect_ratio"
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "get various configuration",
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	getCmd.AddCommand(aspectratio.Get)
	ConfigCMD.AddCommand(getCmd)
}
