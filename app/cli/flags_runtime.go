package cli

func init() {
	rootCmd.PersistentFlags().CountP("verbose", "v", "set verbose level. Set once to print debug level, repeat to print everything")
	rootCmd.Flags().Bool("hold-on-exit", false, "prevent exiting after running a command")
}
