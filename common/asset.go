package common

import (
	"runtime"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

type Asset struct {
	Name      string
	URL       string
	Size      int
	CreatedAt time.Time
	Logger    *logrus.Logger
}

func (a *Asset) MatchSystem() bool {
	if !a.MatchOS() {
		return false
	}

	if !a.MatchArch() {
		return false
	}

	return true
}

func (a *Asset) MatchArch() bool {
	if strings.Contains(a.Name, "-"+runtime.GOARCH+"-") {
		return true
	}

	if runtime.GOARCH == "amd64" {
		if strings.Contains(a.Name, "-x86_64-") {
			return true
		}
	}

	return false
}

func (a *Asset) MatchOS() bool {
	switch runtime.GOOS {
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
	return strings.Contains(a.Name, "-linux-")
}

func matchDarwin(a *Asset) bool {
	return strings.Contains(a.Name, "-darwin-")
}

func matchWindows(a *Asset) bool {
	if strings.Contains(a.Name, "-win32-") {
		return true
	}

	if strings.HasSuffix(a.Name, ".exe") {
		return true
	}

	return false
}
