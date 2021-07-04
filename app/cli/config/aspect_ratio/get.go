package aspectratio

import (
	"fmt"

	"github.com/spf13/cobra"
)

var Get = &cobra.Command{
	Use:   "aspect_ratio",
	Short: "get aspect ratio parameters",
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
		// TODO: this should show aspect_ratio status
	},
}

var getAspectRatioHeight = &cobra.Command{
	Use:   "height",
	Short: "get current aspect ratio height",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("height")
	},
}

var getAspectRatioWidth = &cobra.Command{
	Use:   "width",
	Short: "get current aspect ratio width",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("width")
	},
}

var getAspectRatioStatus = &cobra.Command{
	Use:   "status",
	Short: "get current aspect ratio status",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("status")
	},
}

func init() {
	Get.AddCommand(getAspectRatioHeight)
	Get.AddCommand(getAspectRatioWidth)
	Get.AddCommand(getAspectRatioStatus)
}
