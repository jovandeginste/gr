package app

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/jovandeginste/gr/common"
	"github.com/jovandeginste/gr/github"
	"github.com/sirupsen/logrus"
)

type Fetcher struct {
	Host    string
	Project string
	Retry   bool
	Version *common.Version

	app *App
}

func (f *Fetcher) Destination() *common.Destination {
	return f.app.Destination
}

func (f *Fetcher) Logger() *logrus.Logger {
	return f.app.Logger
}

func (f *Fetcher) ParseURL(u string) error {
	if !strings.Contains(u, "://") {
		u = "https://" + u
	}

	parsed, err := url.Parse(u)
	if err != nil {
		return err
	}

	f.detectHost(parsed.Host)
	f.Project = strings.Trim(parsed.Path, "/")

	f.init()

	return nil
}

func (f *Fetcher) detectHost(h string) {
	if h == "github.com" || strings.HasSuffix(h, ".github.com") {
		f.Host = "github"

		return
	}
}

func (f *Fetcher) init() {
	if f.Version == nil {
		f.Version = common.VersionLatestRelease()
	}
}

func (f *Fetcher) Fetch() error {
	f.init()

	var (
		r   *common.Release
		err error
	)

	f.Logger().Infof("Fetching '%s'...", f.Name())

	switch f.Host {
	case github.Name:
		r, err = f.fetchGithub()
	default:
		return fmt.Errorf("%w: '%s'", common.ErrHostNotKnown, f.Host)
	}

	if err != nil {
		return err
	}

	r.Logger = f.Logger()

	if r == nil {
		return common.ErrNoMatchingRelease
	}

	r.PackageName = f.Name()

	if r.Exists(f.Destination()) {
		if !f.Retry {
			return fmt.Errorf("%w: %s/%s", common.ErrAlreadyDownloaded, r.PackageName, r.Version)
		}

		if err = r.Purge(f.Destination()); err != nil {
			return err
		}
	}

	if err = f.Destination().EnsureDirs(); err != nil {
		return err
	}

	a, err := f.findAsset(r)
	if err != nil {
		return err
	}

	if err := a.DownloadTo(f.Destination()); err != nil {
		return err
	}

	if err := r.DownloadSourceArchiveTo(f.Destination()); err != nil {
		return err
	}

	if err := r.Detect(f.Destination()); err != nil {
		return err
	}

	f.Logger().Infof("Successfully fetched '%s'", f.Name())

	return nil
}

func (f *Fetcher) findAsset(r *common.Release) (*common.Asset, error) {
	for _, a := range r.Assets {
		if a.MatchSystem() {
			if f.app.Preferences == nil || f.app.Preferences.MatchAsset(a) {
				a.Logger = f.Logger()
				a.PackageName = f.Name()
				a.Version = r.Version

				return a, nil
			}
		}
	}

	return nil, common.ErrNoMatchingAsset
}

func (f *Fetcher) fetchGithub() (*common.Release, error) {
	return github.Fetch(f.Logger(), f.Project, f.Version)
}

func (f *Fetcher) Name() string {
	return strings.Join([]string{f.Host, strings.ReplaceAll(f.Project, "/", ".")}, ".")
}
