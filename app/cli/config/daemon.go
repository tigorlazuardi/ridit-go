package config

import (
	"fmt"
	"strconv"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
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
		fmt.Println(val)
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
	Use:   "interval",
	Short: "set interval between wallpaper changes",
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

var enableWallpaperChange = &cobra.Command{
	Use:   "enable",
	Short: "enable wallpaper change",
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

var disableWallpaperChange = &cobra.Command{
	Use:   "disable",
	Short: "disable wallpaper change",
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	daemonWallpaperConfig.AddCommand(daemonWallpaperInterval)
	daemonWallpaperConfig.AddCommand(enableWallpaperChange)
	daemonWallpaperConfig.AddCommand(disableWallpaperChange)
	daemonConfigCMD.AddCommand(daemonWallpaperConfig)
	daemonConfigCMD.AddCommand(daemonPortConfig)
	ConfigCMD.AddCommand(daemonConfigCMD)
}
