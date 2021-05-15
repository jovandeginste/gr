package common

import (
	"runtime"
	"time"
)

type Asset struct {
	Name      string
	URL       string
	Size      int
	CreatedAt time.Time
}

func (a *Asset) MatchSystem() bool {
	os := runtime.GOOS
	switch os {
	case "windows":
		return matchWindows(a)
	case "darwin":
		return matchDarwin(a)
	case "linux":
		return matchLinux(a)
	default:
		return false
	}
}

func matchLinux(a *Asset) bool {
	return true
}

func matchDarwin(a *Asset) bool {
	return false
}

func matchWindows(a *Asset) bool {
	return false
}
