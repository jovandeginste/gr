package common

import "time"

type Release struct {
	Name        string
	Version     string
	Assets      []Asset
	ReleaseType ReleaseType
	CreatedAt   time.Time
	PublishedAt time.Time
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
