package github

import (
	"github.com/jovandeginste/gr/common"
	"github.com/sirupsen/logrus"
)

const (
	Name = "github"
)

func releasesURLFor(project string) string {
	return "https://api.github.com/repos/" + project + "/releases"
}

func (r *release) archiveURL() string {
	return r.TarballURL
}

func Fetch(l *logrus.Logger, project string, version *common.Version) (*common.Release, error) {
	url := releasesURLFor(project)

	l.Debugf("URL is: '%s'", url)

	var releases []release

	if err := common.Fetch(l, url, &releases); err != nil {
		return nil, err
	}

	if len(releases) == 0 {
		return nil, common.ErrNoRelease
	}

	for _, r := range releases {
		if c := r.asCommonRelease(); c.Match(version) {
			return &c, nil
		}
	}

	return nil, common.ErrNoMatchingRelease
}

func (r *release) asCommonRelease() common.Release {
	c := common.Release{
		Name:          r.Name,
		Version:       r.TagName,
		ReleaseType:   common.ReleaseTypeRelease,
		CreatedAt:     r.CreatedAt,
		PublishedAt:   r.PublishedAt,
		SourceArchive: r.archiveURL(),
	}

	for _, a := range r.Assets {
		if a.isReady() {
			c.Assets = append(c.Assets, a.asCommonAsset())
		}
	}

	if r.Draft {
		c.ReleaseType = common.ReleaseTypeDraft
	}

	if r.Prerelease {
		c.ReleaseType = common.ReleaseTypePrerelease
	}

	return c
}

func (a *asset) isReady() bool {
	return a.State == "uploaded"
}

func (a *asset) asCommonAsset() *common.Asset {
	return &common.Asset{
		Name:      a.Name,
		URL:       a.BrowserDownloadURL,
		Size:      a.Size,
		CreatedAt: a.CreatedAt,
	}
}
