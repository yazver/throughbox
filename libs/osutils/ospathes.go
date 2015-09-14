package osutils

import (
	"os"
	"runtime"
)

func GetHomeDir() {
	if runtime.GOOS == "windows" {
		return os.Getenv("UserProfile")
	} else {
		return os.Getenv("HOME")
	}
}