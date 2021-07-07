package pkg

import "os"

var termode = false

func init() {
	fileinfo, _ := os.Stdout.Stat()
	termode = (fileinfo.Mode() & os.ModeCharDevice) != 0
}

func IsTerminal() bool {
	return termode
}
