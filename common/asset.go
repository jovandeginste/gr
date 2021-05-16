package common

import (
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

type Asset struct {
	PackageName string
	Version     string
	Name        string
	URL         string
	Size        int
	CreatedAt   time.Time
	Logger      *logrus.Logger
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
		return a.matchWindows()
	case "darwin":
		return a.matchDarwin()
	case "linux":
		return a.matchLinux()
	default:
		return false
	}
}

func (a *Asset) matchLinux() bool {
	return strings.Contains(a.Name, "-linux-")
}

func (a *Asset) matchDarwin() bool {
	return strings.Contains(a.Name, "-darwin-")
}

func (a *Asset) matchWindows() bool {
	if strings.Contains(a.Name, "-win32-") {
		return true
	}

	if strings.HasSuffix(a.Name, ".exe") {
		return true
	}

	return false
}

func (a *Asset) Exists(destination *Destination) bool {
	extractDir := destination.GetPackageDirFor(a)

	_, err := os.Stat(extractDir)

	return !os.IsNotExist(err)
}
