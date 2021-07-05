package subreddit

import "github.com/spf13/cobra"

var SubredditCMD = &cobra.Command{
	Use:   "subreddit",
	Short: "add or delete subreddits",
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}
