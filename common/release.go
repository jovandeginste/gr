package common

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"time"

	"github.com/sirupsen/logrus"
)

type Release struct {
	PackageName   string
	Name          string
	Version       string
	Assets        []*Asset
	ReleaseType   ReleaseType
	CreatedAt     time.Time
	PublishedAt   time.Time
	SourceArchive string
	Logger        *logrus.Logger
}

type ReleaseType string

const (
	ReleaseTypeRelease    ReleaseType = "release"
	ReleaseTypePrerelease ReleaseType = "prerelease"
	ReleaseTypeDraft      ReleaseType = "draft"
)

func (r *Release) Match(v Version) bool {
	if !v.AllowRelease && r.ReleaseType == ReleaseTypeRelease {
		return false
	}

	if !v.AllowPrerelease && r.ReleaseType == ReleaseTypePrerelease {
		return false
	}

	if !v.AllowDraft && r.ReleaseType == ReleaseTypeDraft {
		return false
	}

	return true
}

func (r *Release) DownloadSourceArchiveTo(destination *Destination) error {
	extractDir := destination.GetSourceDirFor(r)
	if _, err := os.Stat(extractDir); os.IsExist(err) {
		return fmt.Errorf("%w: %s/%s", ErrAlreadyDownloaded, r.Name, r.Version)
	}

	root := destination.GetTmpDir()
	if err := ensureDir(root); err != nil {
		return err
	}

	root, err := ioutil.TempDir(root, "gr-")
	if err != nil {
		return err
	}

	defer os.RemoveAll(root)

	if err := ensureDir(extractDir); err != nil {
		return err
	}

	file := path.Join(root, "source.tar.gz")
	if err := Download(r.Logger, r.SourceArchive, file); err != nil {
		return err
	}

	r.Logger.Infof("Unpacking in '%s'...", extractDir)

	if err := unpack(file, extractDir); err != nil {
		return err
	}

	return nil
}
