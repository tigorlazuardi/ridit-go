package cli

func init() {
	rootCmd.Flags().Bool("aspect-ratio-enable", true, "enable image aspect ratio")
	rootCmd.Flags().UintP("aspect-ratio-width", "W", 16, "set width aspect ratio for the image")
	rootCmd.Flags().UintP("aspect-ratio-height", "H", 9, "set height aspect ratio for the image")
	rootCmd.Flags().Float64("aspect-ratio-range", 0.5, "set the valid ratio range for the image")
}
