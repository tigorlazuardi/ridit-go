package cli

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/tigorlazuardi/ridit-go/app/cli/config"
	"github.com/tigorlazuardi/ridit-go/app/cli/subreddit"
	"github.com/tigorlazuardi/ridit-go/pkg"
)

var rootCmd = &cobra.Command{
	Use:   "ridit",
	Short: "reddit image downloader",
	Long:  "A CLI program to download images from reddit",
	Run: func(cmd *cobra.Command, args []string) {
		logrus.WithField("args", args).Info("testing args")
	},
}

func Exec() {
	if err := rootCmd.Execute(); err != nil {
		logrus.WithError(err).Panic("failed to execute root command")
	}
}

func init() {
	cobra.OnInitialize(initConfigurations)
	rootCmd.AddCommand(config.ConfigCMD)
	rootCmd.AddCommand(subreddit.SubredditCMD)
	rootCmd.PersistentFlags().StringP("profile", "p", "main", "sets the profile to use")
	rootCmd.PersistentFlags().CountP("verbose", "v", "set verbose level. Set once to print debug level, repeat to print everything")
}

func initConfigurations() {
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:            os.Getenv("RIDIT_LOCAL_DEVELOPMENT") != "",
		PadLevelText:           true,
		DisableLevelTruncation: true,
		FullTimestamp:          true,
		TimestampFormat:        "Jan 02 15:04:05",
	})
	logrus.AddHook(&pkg.JSONHook{})
}
