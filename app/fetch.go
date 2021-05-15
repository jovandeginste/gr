package app

import (
	"fmt"

	"github.com/jovandeginste/gr/common"
	"github.com/jovandeginste/gr/github"
)

func Fetch(host, org, project string, version common.Version) (*common.Release, error) {
	switch host { // nolint:gocritic
	case github.Name:
		return github.Fetch(org, project, version)
	}

	return nil, fmt.Errorf("%w: '%s'", common.ErrHostNotKnown, host)
}
