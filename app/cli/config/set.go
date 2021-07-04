package config

import (
	"github.com/spf13/cobra"
	aspectratio "github.com/tigorlazuardi/ridit-go/app/cli/config/aspect_ratio"
)

var setCmd = &cobra.Command{
	Use:   "set",
	Short: "set various configuration",
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	setCmd.AddCommand(aspectratio.Set)
	ConfigCMD.AddCommand(setCmd)
}
