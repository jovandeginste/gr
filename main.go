package main

import (
	"github.com/jovandeginste/gr/app"
	"github.com/jovandeginste/gr/common"
)

func main() {
	f := app.Fetcher{
		Host:        "github",
		Org:         "zellij-org",
		Project:     "zellij",
		Version:     common.VersionLatestRelease(),
		Destination: "~/tmp/",
	}

	if err := f.Fetch(); err != nil {
		panic(err)
	}
}
