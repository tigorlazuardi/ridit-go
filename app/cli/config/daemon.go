package config

import (
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/tigorlazuardi/ridit-go/app/config"
	"github.com/tigorlazuardi/ridit-go/app/config/models"
)

var daemonConfigCMD = &cobra.Command{
	Use:     "daemon",
	Aliases: []string{"http", "server"},
	Short:   "sets various daemon / http configuration",
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

var daemonPortConfig = &cobra.Command{
	Use:   "port",
	Short: "set the port to listen to",
	Long:  "Set the port to listen to. Port must be a value between 1 - 65535. Port 1 - 1023 requires root/admin privilege",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			_ = cmd.Help()
			return
		}
		val, err := strconv.ParseUint(args[0], 10, 64)
		if err != nil {
			logrus.WithField("given_value", args[0]).Fatal("failed to parse value to positive integer value")
		}
		if val < 1 || val > 65535 {
			logrus.Fatal("unsupported port value. Port must be a value between 1 - 65535. Port 1 - 1023 requires root/admin privilege")
		}
		err = config.Modify(func(c *models.Config) {
			c.Daemon.Port = uint(val)
		})
		if err != nil {
			logrus.WithError(err).Fatal("failed to modify configuration")
		}
		logrus.Println("daemon port set to ", val)
	},
}

var daemonWallpaperConfig = &cobra.Command{
	Use:   "wallpaper",
	Short: "set wallppaper configurations",
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

var daemonWallpaperInterval = &cobra.Command{
	Use:     "interval",
	Short:   "set interval between wallpaper changes",
	Long:    "set interval between wallpaper changes. see https://golang.org/pkg/time/#ParseDuration for format",
	Example: "ridit config daemon interval 10m",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			_ = cmd.Help()
			return
		}
		entry := logrus.WithField("usage_example", cmd.Example).WithField("given_value", args[0])
		dur, err := time.ParseDuration(args[0])
		if err != nil {
			entry.WithError(err).
				Fatal("failed to parse time format. see https://golang.org/pkg/time/#ParseDuration for format")
		}
		err = config.Modify(func(c *models.Config) {
			c.Daemon.WallpaperInterval = models.Duration{Duration: dur}
		})
		if err != nil {
			entry.
				WithError(err).
				Fatal("failed to modify configuration")
		}
		logrus.Println("wallpaper interval set to ", dur.String())
	},
}

var enableWallpaperChange = &cobra.Command{
	Use:     "enable",
	Short:   "enable wallpaper change",
	Example: "ridit config daemon wallpaper enable",
	Run: func(cmd *cobra.Command, args []string) {
		err := config.Modify(func(c *models.Config) {
			c.Daemon.WallpaperChange = true
		})
		if err != nil {
			logrus.WithField("usage_example", cmd.Example).Fatal(err)
		}
		logrus.Println("enabled wallpaper change")
	},
}

var disableWallpaperChange = &cobra.Command{
	Use:     "disable",
	Short:   "disable wallpaper change",
	Example: "ridit config daemon wallpaper disable",
	Run: func(cmd *cobra.Command, args []string) {
		err := config.Modify(func(c *models.Config) {
			c.Daemon.WallpaperChange = false
		})
		if err != nil {
			logrus.WithField("usage_example", cmd.Example).Fatal(err)
		}
		logrus.Println("disabled wallpaper change")
	},
}

func init() {
	daemonWallpaperConfig.AddCommand(daemonWallpaperInterval)
	daemonWallpaperConfig.AddCommand(enableWallpaperChange)
	daemonWallpaperConfig.AddCommand(disableWallpaperChange)
	daemonConfigCMD.AddCommand(daemonWallpaperConfig)
	daemonConfigCMD.AddCommand(daemonWallpaperInterval)
	daemonConfigCMD.AddCommand(daemonPortConfig)
	ConfigCMD.AddCommand(daemonConfigCMD)
}
