package config

import "github.com/spf13/cobra"

var minimumSize = &cobra.Command{
	Use:   "minimum_size",
	Short: "set various minimum size configuration",
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

var enableMinimumSize = &cobra.Command{
	Use:   "enable",
	Short: "enable checking image minimum size",
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

var disableMinimumSize = &cobra.Command{
	Use:   "disable",
	Short: "disable checking image minimum size",
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

var setMinimumSizeHeight = &cobra.Command{
	Use:   "height",
	Short: "set current height minimum size",
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

var setMinimumSizeWidth = &cobra.Command{
	Use:   "width",
	Short: "set current width minimum size",
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	minimumSize.AddCommand(disableMinimumSize)
	minimumSize.AddCommand(enableMinimumSize)
	minimumSize.AddCommand(setMinimumSizeHeight)
	minimumSize.AddCommand(setMinimumSizeWidth)
	ConfigCMD.AddCommand(minimumSize)
}
