package common

import (
	"os"

	"github.com/mholt/archiver/v3"
	"github.com/mitchellh/go-homedir"
)

func isExec(mode os.FileMode) bool {
	return mode&0o111 != 0
}

func expand(dir string) string {
	dir, _ = homedir.Expand(dir)

	return dir
}

func unpack(file, extractDir string) error {
	return archiver.Unarchive(file, extractDir)
}

func ensureDir(p string) error {
	_, err := os.Stat(p)
	if os.IsNotExist(err) {
		return os.MkdirAll(p, 0o700)
	}

	return err
}
