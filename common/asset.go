package common

import (
	"io/ioutil"
	"os"
	"path"
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
	if strings.Contains(a.Name, runtime.GOARCH) {
		return true
	}

	if runtime.GOARCH == "amd64" {
		if strings.Contains(a.Name, "x86_64") {
			return true
		}

		if strings.Contains(a.Name, "AppImage") {
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
	if strings.Contains(a.Name, "AppImage") {
		return true
	}

	return strings.Contains(a.Name, "linux")
}

func (a *Asset) matchDarwin() bool {
	return strings.Contains(a.Name, "darwin")
}

func (a *Asset) matchWindows() bool {
	if strings.Contains(a.Name, "win32") {
		return true
	}

	if strings.HasSuffix(a.Name, ".exe") {
		return true
	}

	return false
}

func (a *Asset) DownloadTo(destination *Destination) error {
	root := destination.GetTmpDir()
	if err := ensureDir(root); err != nil {
		return err
	}

	root, err := ioutil.TempDir(root, "gr-")
	if err != nil {
		return err
	}

	defer os.RemoveAll(root)

	extractDir := destination.GetReleaseDirFor(a)
	if err := ensureDir(extractDir); err != nil {
		return err
	}

	file := path.Join(root, a.Name)
	if err := Download(a.Logger, a.URL, file); err != nil {
		return err
	}

	if isArchive(file) {
		a.Logger.Debugf("Unpacking in '%s'...", extractDir)

		return unpack(file, extractDir)
	}

	if isAppImage(file) {
		a.Logger.Debugf("Copying to '%s'...", extractDir)

		return moveAppImage(file, extractDir)
	}

	return nil
}
