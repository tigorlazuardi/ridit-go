package cli

import (
	"os/user"
	"path"
	"time"

	"github.com/kirsle/configdir"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var daemonCMD = &cobra.Command{
	Use:     "daemon",
	Aliases: []string{"http", "serve"},
	Short:   "runs ridit as daemon service in port 10101",
	Long:    "runs ridit as daemon service in port 10101",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() {
	configDir := configdir.LocalConfig("ridit")
	daemonCMD.Flags().Bool("change-wallpaper", false, "change wallpaper periodically")
	daemonCMD.Flags().Duration("change-wallpaper-interval", time.Minute*10, "set the time interval between wallpaper changing. see https://golang.org/pkg/time/#ParseDuration for format")
	daemonCMD.Flags().String("db", configDir, "sets the db location")
	usr, err := user.Current()
	if err != nil {
		logrus.WithError(err).Panic("cannot get current user information")
	}
	dp := path.Join(usr.HomeDir, "Pictures", "ridit")
	daemonCMD.Flags().String("download-path", dp, "set the download path")

	rootCmd.AddCommand(daemonCMD)
}
