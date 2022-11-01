package common

import (
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

func (r *Release) Match(v *Version) bool {
	if v == nil {
		return true
	}

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
	root := destination.GetTmpDir()
	if err := ensureDir(root); err != nil {
		return err
	}

	root, err := ioutil.TempDir(root, "gr-")
	if err != nil {
		return err
	}

	defer os.RemoveAll(root)

	extractDir := destination.GetSourceDirFor(r)
	if err := ensureDir(extractDir); err != nil {
		return err
	}

	file := path.Join(root, "source.tar.gz")
	if err := Download(r.Logger, r.SourceArchive, file); err != nil {
		return err
	}

	r.Logger.Debugf("Unpacking in '%s'...", extractDir)

	if err := unpack(file, extractDir); err != nil {
		return err
	}

	return nil
}

func (r *Release) Exists(destination *Destination) bool {
	extractDir := destination.GetPackageDirFor(r.PackageName, r.Version)

	_, err := os.Stat(extractDir)

	return err == nil || os.IsExist(err)
}

func (r *Release) Purge(destination *Destination) error {
	extractDir := destination.GetPackageDirFor(r.PackageName, r.Version)

	r.Logger.Infof("Removing package dir '%s'...", extractDir)

	err := os.RemoveAll(extractDir)

	return err
}
