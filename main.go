package main

import (
	"fmt"

	"github.com/jovandeginste/gr/app"
	"github.com/jovandeginste/gr/common"
)

func main() {
	host := "github"
	org := "zellij-org"
	project := "zellij"
	version := common.VersionLatestRelease()

	r, err := app.Fetch(host, org, project, version)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%#v\n", r)
}
