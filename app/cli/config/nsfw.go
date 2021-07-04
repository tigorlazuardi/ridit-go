package config

import "github.com/spf13/cobra"

var nsfwCMD = &cobra.Command{
	Use:     "nsfw",
	Aliases: []string{"adult", "over_18"},
	Short:   "toggle nsfw images",
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

var enableNSFW = &cobra.Command{
	Use:     "enable",
	Aliases: []string{"enabled"},
	Short:   "enable nsfw image",
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

var disableNSFW = &cobra.Command{
	Use:     "disable",
	Aliases: []string{"enabled"},
	Short:   "disable nsfw images",
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	nsfwCMD.AddCommand(enableNSFW)
	nsfwCMD.AddCommand(disableNSFW)
	ConfigCMD.AddCommand(nsfwCMD)
}
