package config

import "github.com/spf13/cobra"

var downloadConfig = &cobra.Command{
	Use:   "download",
	Short: "sets various download configuration",
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

var pathConfig = &cobra.Command{
	Use:     "path",
	Short:   "sets download path",
	Example: "ridit config download path '~/Pictures'",
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

var connectTimeout = &cobra.Command{
	Use:     "connect_timeout",
	Short:   "Duration to wait establishing connection to source download",
	Example: "ridit config download connect_timeout 5s",
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

var downloadTimeout = &cobra.Command{
	Use:     "timeout",
	Short:   "Duration to wait for download to finish. Increase if you have slow connection",
	Example: "ridit config download timeout 30s",
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	downloadConfig.AddCommand(connectTimeout)
	downloadConfig.AddCommand(downloadTimeout)
	downloadConfig.AddCommand(pathConfig)
	ConfigCMD.AddCommand(downloadConfig)
}
