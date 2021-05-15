package common

import "golang.org/x/mod/semver"

type (
	VersionMatcher func(string, string) int

	version struct {
		AllowRelease, AllowDraft, AllowPrerelease bool
		Version                                   string
		Matcher                                   VersionMatcher
	}
	Version struct {
		*version
	}
)

func VersionLatestReleaseAfter(v string) Version {
	return Version{
		version: &version{
			Version:      v,
			AllowRelease: true,
			Matcher:      semver.Compare,
		},
	}
}

func VersionLatestRelease() Version {
	return Version{
		version: &version{
			AllowRelease: true,
		},
	}
}

func (v *Version) Match(other Release) int {
	if v.Matcher() == nil {
		return 0
	}

	return v.Matcher()(v.Version, other.Version)
}

func (v *Version) Matcher() VersionMatcher {
	return v.version.Matcher
}
