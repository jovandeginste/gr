package app

import (
	"fmt"
	"path"

	"github.com/jovandeginste/gr/common"
	"github.com/jovandeginste/gr/github"
	"github.com/sirupsen/logrus"
)

type Fetcher struct {
	Destination *common.Destination
	Host        string
	Org         string
	Project     string
	Version     common.Version
	Preferences *common.Preferences
	Logger      *logrus.Logger
}

func (f *Fetcher) init() {
	if f.Logger == nil {
		f.Logger = logrus.New()
	}
}

func (f *Fetcher) Fetch() error {
	f.init()

	var (
		r   *common.Release
		err error
	)

	switch f.Host {
	case github.Name:
		r, err = f.fetchGithub()
	default:
		return fmt.Errorf("%w: '%s'", common.ErrHostNotKnown, f.Host)
	}

	r.Logger = f.Logger

	if err != nil {
		return err
	}

	if r == nil {
		return common.ErrNoMatchingRelease
	}

	r.PackageName = f.Name()

	if r.Exists(f.Destination) {
		return fmt.Errorf("%w: %s/%s", common.ErrAlreadyDownloaded, r.PackageName, r.Version)
	}

	if err := f.Destination.EnsureDirs(); err != nil {
		return err
	}

	a, err := f.findAsset(r)
	if err != nil {
		return err
	}

	if err := a.DownloadTo(f.Destination); err != nil {
		return err
	}

	if err := r.DownloadSourceArchiveTo(f.Destination); err != nil {
		return err
	}

	return r.Detect(f.Destination)
}

func (f *Fetcher) findAsset(r *common.Release) (*common.Asset, error) {
	for _, a := range r.Assets {
		if a.MatchSystem() {
			if f.Preferences == nil || f.Preferences.MatchAsset(a) {
				a.Logger = f.Logger
				a.PackageName = f.Name()
				a.Version = r.Version

				return a, nil
			}
		}
	}

	return nil, common.ErrNoMatchingAsset
}

func (f *Fetcher) fetchGithub() (*common.Release, error) {
	return github.Fetch(f.Logger, f.Org, f.Project, f.Version)
}

func (f *Fetcher) Name() string {
	return path.Join(f.Host, f.Org, f.Project)
}
