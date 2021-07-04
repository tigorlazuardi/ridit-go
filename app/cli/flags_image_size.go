package cli

func init() {
	rootCmd.PersistentFlags().Bool("minimum-size-enable", true, "enable minimum image size checking")
	rootCmd.PersistentFlags().Uint("minimum-size-width", 1920, "sets minimum width for images. set to 0 to allow all width")
	rootCmd.PersistentFlags().Uint("minimum-size-height", 1080, "sets minimum height for images. set to 0 to allow all height")
}
