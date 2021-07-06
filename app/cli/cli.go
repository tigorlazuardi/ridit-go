package cli

import (
	"context"
	"net/http"
	"os"
	"path/filepath"

	"github.com/kirsle/configdir"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tigorlazuardi/ridit-go/app/cli/config"
	"github.com/tigorlazuardi/ridit-go/app/cli/subreddit"
	configapi "github.com/tigorlazuardi/ridit-go/app/config"
	"github.com/tigorlazuardi/ridit-go/app/reddit"
	"github.com/tigorlazuardi/ridit-go/pkg"
)

var rootCmd = &cobra.Command{
	Use:   "ridit",
	Short: "reddit image downloader",
	Long:  "A CLI program to download images from reddit",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		entry := pkg.EntryFromContext(ctx)
		config, err := configapi.Load()
		if err != nil {
			entry.WithError(err).Fatal("failed to read config file")
		}

		repository := reddit.NewRepository(http.DefaultClient, config)
		for fChan := range repository.GetListing(ctx) {
			if fChan.Err != nil {
				err = fChan.Err
				entry.WithError(err).Error(err)
			} else {
				entry.WithField("data", fChan.Downloads).Debug("data")
			}
		}
	},
}

func Exec() {
	ctx := pkg.ContextWithNewEntry(context.Background())
	if err := rootCmd.ExecuteContext(ctx); err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfigurations)
	rootCmd.AddCommand(config.ConfigCMD)
	rootCmd.AddCommand(subreddit.SubredditCMD)
	rootCmd.PersistentFlags().StringP("profile", "p", "main", "sets the profile to use")
	rootCmd.PersistentFlags().CountP("verbose", "v", "set verbose level. Set once to print debug level, repeat to print everything")
	rootCmd.PersistentFlags().BoolP("machine", "m", false, "runs the cli program with assumption machine is executing the application. logs will now appear in json format")
}

func initConfigurations() {
	dev := os.Getenv("RIDIT_LOCAL_DEVELOPMENT") != ""
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:            dev,
		PadLevelText:           true,
		DisableLevelTruncation: true,
		FullTimestamp:          true,
		TimestampFormat:        "Jan 02 15:04:05",
	})
	if dev {
		logrus.SetLevel(logrus.DebugLevel)
	}
	if machine, _ := rootCmd.PersistentFlags().GetBool("machine"); machine {
		logrus.SetFormatter(&logrus.JSONFormatter{})
	}
	logrus.AddHook(&pkg.JSONHook{})
	logrus.AddHook(&pkg.FrameHook{Disabled: !dev})
	prof, _ := rootCmd.Flags().GetString("profile")
	dir := configdir.LocalConfig("ridit")
	err := os.MkdirAll(dir, 0777)
	if err != nil {
		logrus.WithError(err).Fatal("failed to create configuration folder on ", dir)
	}
	err = viper.BindPFlags(rootCmd.PersistentFlags())
	if err != nil {
		logrus.WithError(err).Fatal("failed to bind flags from cobra")
	}
	viper.Set("configfile", filepath.Join(dir, prof+".toml"))

	file, created, err := configapi.LoadConfigFile()
	if err != nil {
		logrus.WithError(err).Fatal("failed to create config file")
	}
	defer file.Close()

	if created {
		logrus.WithField("location", viper.GetString("configfile")).Info("config file created")
		os.Exit(0)
	}

	viper.SetConfigType("toml")
	err = viper.ReadConfig(file)
	if err != nil {
		logrus.WithError(err).Fatal("failed to read config file")
	}
}
