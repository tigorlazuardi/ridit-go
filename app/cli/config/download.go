package config

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tigorlazuardi/ridit/app/config"
	"github.com/tigorlazuardi/ridit/app/config/models"
	"github.com/tigorlazuardi/ridit/pkg"
)

var downloadConfig = &cobra.Command{
	Use:   "download",
	Short: "sets various download configuration",
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

var pathConfig = &cobra.Command{
	Use:     "path",
	Short:   "sets download path",
	Example: "ridit config download path '~/Pictures'",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			_ = cmd.Help()
			os.Exit(1)
		}
		entry := pkg.EntryFromContext(cmd.Context())
		err := config.Modify(viper.GetString("profile"), func(c *models.Config) {
			c.Download.Path = args[0]
		})
		if err != nil {
			entry.Fatal(err)
		}
		fmt.Println("download path set to ", args[0])
	},
}

var connectTimeout = &cobra.Command{
	Use:     "connect_timeout",
	Short:   "Duration to wait establishing connection to source download",
	Example: "ridit config download connect_timeout 5s",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			_ = cmd.Help()
			os.Exit(1)
		}

		entry := pkg.EntryFromContext(cmd.Context())

		t, err := time.ParseDuration(args[0])
		if err != nil {
			entry.Fatal(err)
		}
		err = config.Modify(viper.GetString("profile"), func(c *models.Config) {
			c.Download.ConnectTimeout.Duration = t
		})
		if err != nil {
			entry.Fatal(err)
		}

		fmt.Println("download connection timeout set to ", args[0])
	},
}

var downloadTimeout = &cobra.Command{
	Use:     "timeout",
	Short:   "Duration to wait for download to finish. Increase if you have slow connection",
	Example: "ridit config download timeout 30s",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			_ = cmd.Help()
			os.Exit(1)
		}

		entry := pkg.EntryFromContext(cmd.Context())

		t, err := time.ParseDuration(args[0])
		if err != nil {
			entry.Fatal(err)
		}
		err = config.Modify(viper.GetString("profile"), func(c *models.Config) {
			c.Download.Timeout.Duration = t
		})
		if err != nil {
			entry.Fatal(err)
		}

		fmt.Println("download timeout set to ", args[0])
	},
}

func init() {
	downloadConfig.AddCommand(connectTimeout)
	downloadConfig.AddCommand(downloadTimeout)
	downloadConfig.AddCommand(pathConfig)
	ConfigCMD.AddCommand(downloadConfig)
}
