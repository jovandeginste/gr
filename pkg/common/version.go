package common

import "golang.org/x/mod/semver"

type (
	VersionMatcher func(string, string) int

	Version struct {
		AllowRelease, AllowDraft, AllowPrerelease bool
		Version                                   string
		Matcher                                   VersionMatcher
	}
)

func VersionLatestReleaseAfter(v string) *Version {
	return &Version{
		Version:      v,
		AllowRelease: true,
		Matcher:      semver.Compare,
	}
}

func VersionLatestRelease() *Version {
	return &Version{
		AllowRelease: true,
	}
}

func (v *Version) Match(other Release) int {
	if v.Matcher == nil {
		return 0
	}

	return v.Matcher(v.Version, other.Version)
}
