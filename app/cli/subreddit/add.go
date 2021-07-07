package subreddit

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tigorlazuardi/ridit/app/config"
	"github.com/tigorlazuardi/ridit/app/config/models"
	"github.com/tigorlazuardi/ridit/app/config/models/sort"
	"github.com/tigorlazuardi/ridit/app/reddit"
	"github.com/tigorlazuardi/ridit/pkg"
)

var addSubreddit = &cobra.Command{
	Use:     "add",
	Aliases: []string{"insert", "install"},
	Short:   "adds subreddit to list. case-insensitive",
	Long:    "adds subreddit to list. case-insensitive. Replaces old value if exist",
	Example: "ridit subreddit add animewallpaper --sort=new --nsfw=false",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		profile := viper.GetString("profile")
		if len(args) == 0 {
			_ = cmd.Help()
			os.Exit(1)
		}
		entry := pkg.EntryFromContext(ctx)
		rawSort, _ := cmd.Flags().GetString("sort")
		sort, err := sort.Parse(rawSort)
		if err != nil {
			entry.WithFields(logrus.Fields{
				"usage_example": cmd.Example,
				"given_value":   rawSort,
			}).Fatal(err)
		}
		nsfw, _ := cmd.Flags().GetBool("nsfw")
		subs := make(map[string]models.Subreddit)
		if noCheck, _ := cmd.Flags().GetBool("no-check"); !noCheck {
			for check := range reddit.CheckSubredditExist(nil, ctx, args) {
				entry.WithField("check", check).Debug("check value")
				if check.Err != nil {
					entry.WithError(check.Err).Error(check.Err)
					continue
				}
				if check.Exist {
					subs[check.Name] = models.Subreddit{
						Sort: sort,
						NSFW: nsfw,
					}
				}
			}
		} else {
			for _, v := range args {
				subs[v] = models.Subreddit{
					Sort: sort,
					NSFW: nsfw,
				}
			}
		}
		err = config.Modify(profile, func(c *models.Config) {
			for k, v := range subs {
				c.Subreddits[k] = v
			}
		})
		if err != nil {
			entry.WithField("usage_example", cmd.Example).Fatal(err)
		}
		entry.WithField("subreddits", args).Info("subreddit(s) added")
	},
}

func init() {
	addSubreddit.Flags().Bool("no-check", false, "skip validation checking")
	addSubreddit.Flags().StringP("sort", "s", "new", "set sort value")
	addSubreddit.Flags().BoolP("nsfw", "n", true, "set enable downloading nsfw images to subreddit. has no effect if global nsfw is disabled")
	SubredditCMD.AddCommand(addSubreddit)
}
