package cli

import "os"

func init() {
	rootCmd.Flags().Bool("nsfw", true, "enable / disable nsfw images")
	rootCmd.Flags().StringP("sort", "s", "new", "sort subreddit. valid values: 'new', 'hot', 'controversial'")
	rootCmd.Flags().BoolP("enable-symlink", "l", false, "enable symlinking images to 'symlink-path'")
	cwd, _ := os.Getwd()
	rootCmd.Flags().String("symlink-path", cwd, "symlink target location. default's to current working directory (cwd)")
}
