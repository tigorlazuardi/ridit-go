package cli

import "os"

func init() {
	rootCmd.PersistentFlags().Bool("nsfw", true, "enable / disable nsfw images")
	rootCmd.PersistentFlags().StringP("sort", "s", "new", "sort subreddit. valid values: 'new', 'hot', 'controversial'")
	rootCmd.PersistentFlags().BoolP("enable-symlink", "l", false, "enable symlinking images to 'symlink-path'")
	cwd, _ := os.Getwd()
	rootCmd.PersistentFlags().String("symlink-path", cwd, "symlink target location. default's to current working directory (cwd)")
}
