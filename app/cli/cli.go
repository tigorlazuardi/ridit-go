package cli

import (
	"fmt"
	"os"
	"path"

	"github.com/kirsle/configdir"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tigorlazuardi/ridit-go/app/cli/config"
	"github.com/tigorlazuardi/ridit-go/app/cli/subreddit"
	configapi "github.com/tigorlazuardi/ridit-go/app/config"
	"github.com/tigorlazuardi/ridit-go/pkg"
)

var rootCmd = &cobra.Command{
	Use:   "ridit",
	Short: "reddit image downloader",
	Long:  "A CLI program to download images from reddit",
	Run: func(cmd *cobra.Command, args []string) {
		config := configapi.Load()

		fmt.Println(config.Download.Path)
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
	prof, _ := rootCmd.Flags().GetString("profile")
	dir := configdir.LocalConfig("ridit", prof)
	err := os.MkdirAll(dir, os.ModeDir)
	if err != nil {
		logrus.WithError(err).Fatal("failed to create configuration folder on ", dir)
	}
	err = viper.BindPFlags(rootCmd.PersistentFlags())
	if err != nil {
		logrus.WithError(err).Fatal("failed to bind flags from cobra")
	}
	viper.Set("configfile", path.Join(dir, configapi.Filename))

	file, created := configapi.LoadConfigFile()
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
