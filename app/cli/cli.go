package cli

import (
	"context"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/ahmetalpbalkan/go-cursor"
	"github.com/kirsle/configdir"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tigorlazuardi/ridit/app/cli/config"
	"github.com/tigorlazuardi/ridit/app/cli/subreddit"
	configapi "github.com/tigorlazuardi/ridit/app/config"
	"github.com/tigorlazuardi/ridit/app/reddit"
	"github.com/tigorlazuardi/ridit/pkg"
)

var rootCmd = &cobra.Command{
	Use:   "ridit",
	Short: "reddit image downloader",
	Long:  "A CLI program to download images from reddit",
	Run: func(cmd *cobra.Command, _ []string) {
		cursor.Hide()
		defer func() {
			cursor.Show()
			if pkg.IsTerminal() {
				time.Sleep(100 * time.Millisecond)
			}
		}()
		ctx := cmd.Context()
		sig := pkg.RegisterInterrupt()
		ctx = pkg.ContextWithInterrupt(ctx, sig)
		entry := pkg.EntryFromContext(ctx)
		profile := viper.GetString("profile")
		config, err := configapi.Load(profile)
		if err != nil {
			entry.WithError(err).Fatal("failed to read config file")
		}

		repository := reddit.NewRepository(http.DefaultClient, config)
		for fChan := range repository.Fetch(ctx) {
			if fChan.Err != nil {
				err = fChan.Err
				entry.WithError(err).Error(err)
			} else {
				entry.WithField("data", fChan.Download).Trace("operation done")
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
	rootCmd.PersistentFlags().Uint("concurrency", 8, "set maximum number of concurrency when downloading images")
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
	if !pkg.IsTerminal() {
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

	file, created, err := configapi.LoadConfigFile(prof)
	if err != nil {
		logrus.WithError(err).Fatal("failed to create config file")
	}
	if created {
		file.Close()
		logrus.WithField("location", viper.GetString("configfile")).Info("config file created")
		os.Exit(0)
	}

	viper.SetConfigType("toml")
	err = viper.ReadConfig(file)
	if err != nil {
		file.Close()
		logrus.WithError(err).Fatal("failed to read config file")
	}
	file.Close()
}
