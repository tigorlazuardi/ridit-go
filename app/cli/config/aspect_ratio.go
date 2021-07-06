package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/tigorlazuardi/ridit-go/app/config"
	"github.com/tigorlazuardi/ridit-go/app/config/models"
)

var aspectRatioCmd = &cobra.Command{
	Use:     "aspect_ratio",
	Short:   "set aspect_ratio parameters",
	Example: "ridit config aspect_ratio enable",
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
		os.Exit(1)
	},
}

var enableAspectRatio = &cobra.Command{
	Use:     "enable",
	Aliases: []string{"enabled", "true", "True"},
	Short:   "enable image aspect ratio check",
	Example: "ridit config aspect_ratio enable",
	Run: func(cmd *cobra.Command, args []string) {
		err := config.Modify(func(c *models.Config) {
			c.AspectRatio.Enabled = true
		})
		if err != nil {
			logrus.WithError(err).WithField("usage_example", cmd.Example).Fatal(err)
		}
	},
}

var disableAspectRatio = &cobra.Command{
	Use:     "disable",
	Aliases: []string{"disabled", "false", "False"},
	Short:   "disable image aspect ratio check",
	Example: "ridit config aspect_ratio disable",
	Run: func(cmd *cobra.Command, args []string) {
		err := config.Modify(func(c *models.Config) {
			c.AspectRatio.Enabled = false
		})
		if err != nil {
			logrus.WithError(err).WithField("usage_example", cmd.Example).Fatal(err)
		}
	},
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
