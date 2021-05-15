package common

import "errors"

var (
	ErrHostNotKnown      = errors.New("host not known")
	ErrNoSuchProject     = errors.New("no such project found")
	ErrNoRelease         = errors.New("no releases found")
	ErrNoMatchingRelease = errors.New("no release found matching constraints")
)
