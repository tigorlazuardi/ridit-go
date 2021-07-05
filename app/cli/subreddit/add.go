package subreddit

import "github.com/spf13/cobra"

var addSubreddit = &cobra.Command{
	Use:     "add",
	Aliases: []string{"insert", "install"},
	Short:   "adds subreddit to list. case-insensitive",
	Long:    "adds subreddit to list. case-insensitive. Replaces old value if exist",
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	addSubreddit.Flags().StringP("sort", "s", "new", "set sort value")
	addSubreddit.Flags().BoolP("nsfw", "n", true, "set enable downloading nsfw images to subreddit. has no effect if global nsfw is disabled")
	SubredditCMD.AddCommand(addSubreddit)
}
