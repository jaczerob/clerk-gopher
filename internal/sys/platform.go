package sys

import (
	"runtime"

	log "github.com/sirupsen/logrus"
)

func GetPlatform() (platform string) {
	platform = runtime.GOOS

	if platform == "windows" {
		arch := runtime.GOARCH
		if arch[len(arch)-2:] == "64" {
			platform = "win64"
		} else {
			platform = "win32"
		}
	}

	log.WithField("platform", platform).Trace("got platform")
	return
}
