package subreddit

import "github.com/spf13/cobra"

var deleteSubreddit = &cobra.Command{
	Use:     "delete",
	Aliases: []string{"remove", "rm"},
	Short:   "delete subreddit from list if exist",
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	SubredditCMD.AddCommand(deleteSubreddit)
}
