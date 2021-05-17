package main

import (
	"github.com/jovandeginste/gr/app"
	"github.com/jovandeginste/gr/common"
)

func main() {
	f := app.Fetcher{
		Host:        "github",
		Org:         "BurntSushi",
		Project:     "ripgrep",
		Version:     common.VersionLatestRelease(),
		Destination: common.NewDestination("~/tmp/gr"),
	}

	if err := f.Fetch(); err != nil {
		panic(err)
	}

	f = app.Fetcher{
		Host:        "github",
		Org:         "dandavison",
		Project:     "delta",
		Version:     common.VersionLatestRelease(),
		Destination: common.NewDestination("~/tmp/gr"),
	}

	if err := f.Fetch(); err != nil {
		panic(err)
	}
}
