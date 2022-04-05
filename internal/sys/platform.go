package sys

import (
	"os"
	"path"
	"runtime"
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

	return
}

func GetDirectory() (dir string, err error) {
	switch GetPlatform() {
	case "win32":
		dir = "C:\\Program Files\\Toontown Rewritten"
	case "win64":
		dir = "C:\\Program Files (x86)\\Toontown Rewritten"
	case "darwin":
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}

		dir = path.Join(home, "/Library/Application Support/Toontown Rewritten")

	default:
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}

		dir = path.Join(home, "/Toontown Rewritten")
	}

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0700)
		return dir, err
	}

	return dir, err
}

func GetExecutable() string {
	switch GetPlatform() {
	case "win32":
		return "TTREngine.exe"
	case "win64":
		return "TTREngine64.exe"
	case "darwin":
		return "Toontown Rewritten"
	default:
		return "TTREngine"
	}
}
