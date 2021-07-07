package subreddit

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tigorlazuardi/ridit/app/config"
	"github.com/tigorlazuardi/ridit/app/config/models"
	"github.com/tigorlazuardi/ridit/pkg"
)

var deleteSubreddit = &cobra.Command{
	Use:     "delete",
	Aliases: []string{"remove", "rm"},
	Short:   "delete subreddit from list if exist",
	Run: func(cmd *cobra.Command, args []string) {
		entry := pkg.EntryFromContext(cmd.Context())
		if len(args) == 0 {
			_ = cmd.Help()
			os.Exit(1)
		}
		err := config.Modify(viper.GetString("profile"), func(c *models.Config) {
			for _, v := range args {
				delete(c.Subreddits, v)
			}
		})
		if err != nil {
			entry.Fatal(err)
		}
		entry.WithField("subreddits", args).Info("subreddits deleted")
	},
}

func init() {
	SubredditCMD.AddCommand(deleteSubreddit)
}
