package cli

import "github.com/kirsle/configdir"

func init() {
	configPath := configdir.LocalConfig("ridit")
	rootCmd.Flags().StringP("config", "c", configPath, "reads config from custom location")
	rootCmd.Flags().Bool("print-config", false, "prints default config to stdout")
}
