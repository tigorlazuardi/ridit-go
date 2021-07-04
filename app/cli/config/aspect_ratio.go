package config

import (
	"fmt"
	"strconv"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var aspectRatioCmd = &cobra.Command{
	Use:   "aspect_ratio",
	Short: "set aspect_ratio parameters",
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

var enableAspectRatio = &cobra.Command{
	Use:     "enable",
	Aliases: []string{"enabled", "true", "True"},
	Short:   "enable image aspect ratio check",
	Run:     func(cmd *cobra.Command, args []string) {},
}

var disableAspectRatio = &cobra.Command{
	Use:     "disable",
	Aliases: []string{"disabled", "false", "False"},
	Short:   "disable image aspect ratio check",
	Run:     func(cmd *cobra.Command, args []string) {},
}

var setAspectRatioHeight = &cobra.Command{
	Use:     "height",
	Short:   "set height aspect ratio",
	Example: "ridit config set aspect_ratio height 1080",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			_ = cmd.Help()
			return
		}
		val, err := strconv.ParseUint(args[0], 10, 64)
		if err != nil {
			logrus.WithField("given_value", args[0]).Fatal("failed to parse value to positive integer value")
		}
		fmt.Println(val)
		// TODO: implement set aspect ratio
	},
}

var setAspectRatioWidth = &cobra.Command{
	Use:     "width",
	Short:   "set width aspect ratio",
	Example: "ridit config set aspect_ratio width 1920",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			_ = cmd.Help()
			return
		}
		val, err := strconv.ParseUint(args[0], 10, 64)
		if err != nil {
			logrus.WithField("given_value", args[0]).Fatal("failed to parse value to positive integer value")
		}
		fmt.Println(val)
		// TODO: implement set aspect ratio
	},
}

func init() {
	aspectRatioCmd.AddCommand(enableAspectRatio)
	aspectRatioCmd.AddCommand(disableAspectRatio)
	aspectRatioCmd.AddCommand(setAspectRatioHeight)
	aspectRatioCmd.AddCommand(setAspectRatioWidth)
	ConfigCMD.AddCommand(aspectRatioCmd)
}
