package common

import (
	"strings"
)

type (
	Preferences struct {
		glibc *bool
		musl  *bool
	}
)

func NoPreferences() *Preferences {
	return &Preferences{}
}

func (p *Preferences) SetGlibc(v bool) {
	p.glibc = &v
}

func (p *Preferences) SetMusl(v bool) {
	p.musl = &v
}

func (p *Preferences) MatchAsset(a *Asset) bool {
	return p.matchLibC(a)
}

func (p *Preferences) matchLibC(a *Asset) bool {
	if p.glibc != nil {
		if *p.glibc != strings.Contains(a.Name, "glibc") {
			return false
		}
	}

	if p.musl != nil {
		if *p.musl != strings.Contains(a.Name, "musl") {
			return false
		}
	}

	return true
}
