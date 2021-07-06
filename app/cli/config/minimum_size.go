package config

import (
	"os"
	"strconv"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/tigorlazuardi/ridit-go/app/config"
	"github.com/tigorlazuardi/ridit-go/app/config/models"
)

var minimumSize = &cobra.Command{
	Use:     "minimum_size",
	Short:   "set various minimum size configuration",
	Example: "ridit config minimum_size enable",
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
		os.Exit(1)
	},
}

var enableMinimumSize = &cobra.Command{
	Use:     "enable",
	Short:   "enable checking image minimum size",
	Aliases: []string{"enabled", "true", "True"},
	Example: "ridit config minimum_size enable",
	Run: func(cmd *cobra.Command, args []string) {
		err := config.Modify(func(c *models.Config) {
			c.MinimumSize.Enabled = true
		})
		if err != nil {
			logrus.WithField("usage_example", cmd.Example).Fatal(err)
		}
		logrus.Println("enabled minimum_size check")
	},
}

var disableMinimumSize = &cobra.Command{
	Use:     "disable",
	Short:   "disable checking image minimum size",
	Aliases: []string{"disabled", "false", "False"},
	Example: "ridit config minimum_size disable",
	Run: func(cmd *cobra.Command, args []string) {
		err := config.Modify(func(c *models.Config) {
			c.MinimumSize.Enabled = false
		})
		if err != nil {
			logrus.WithField("usage_example", cmd.Example).Fatal(err)
		}
		logrus.Println("disabled minimum_size check")
	},
}

var setMinimumSizeHeight = &cobra.Command{
	Use:     "height",
	Short:   "set current height minimum size",
	Example: "ridit config minimum_size height 1080",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			_ = cmd.Help()
			return
		}
		entry := logrus.WithField("usage_example", cmd.Example)
		val, err := strconv.ParseUint(args[0], 10, 64)
		if err != nil {
			entry.WithError(err).Fatal("failed to parse height to positive integer value")
		}
		err = config.Modify(func(c *models.Config) {
			c.MinimumSize.Height = uint(val)
		})
		if err != nil {
			entry.WithError(err).Fatal("failed to modify configuration")
		}
		logrus.Println("minimum height size set to", val)
	},
}

var setMinimumSizeWidth = &cobra.Command{
	Use:     "width",
	Short:   "set current width minimum size",
	Example: "ridit config minimum_size width 1920",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			_ = cmd.Help()
			return
		}
		entry := logrus.WithField("usage_example", cmd.Example)
		val, err := strconv.ParseUint(args[0], 10, 64)
		if err != nil {
			entry.WithError(err).Fatal("failed to parse width to positive integer value")
		}
		err = config.Modify(func(c *models.Config) {
			c.MinimumSize.Width = uint(val)
		})
		if err != nil {
			entry.WithError(err).Fatal("failed to modify configuration")
		}
		logrus.Println("minimum width size set to", val)
	},
}

func init() {
	minimumSize.AddCommand(disableMinimumSize)
	minimumSize.AddCommand(enableMinimumSize)
	minimumSize.AddCommand(setMinimumSizeHeight)
	minimumSize.AddCommand(setMinimumSizeWidth)
	ConfigCMD.AddCommand(minimumSize)
}
