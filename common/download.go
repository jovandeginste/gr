package common

import (
	"io/ioutil"
	"os"
	"path"

	archiver "github.com/mholt/archiver/v3"
)

func (a *Asset) DownloadTo(destination string) error {
	root, e := ioutil.TempDir(os.TempDir(), "gr-")
	if e != nil {
		return e
	}

	defer os.RemoveAll(root)

	downloadDir := path.Join(root, "download")
	if err := os.Mkdir(downloadDir, 0o700); err != nil {
		return err
	}

	extractDir := path.Join(root, "extract")
	if err := os.Mkdir(extractDir, 0o700); err != nil {
		return err
	}

	file := path.Join(downloadDir, a.Name)
	if err := Download(a.Logger, a.URL, file); err != nil {
		return err
	}

	a.Logger.Infof("Unpacking '%s' to '%s'...", file, extractDir)
	if err := unpack(file, extractDir); err != nil {
		return err
	}

	return nil
}

func unpack(file, extractDir string) error {
	return archiver.Unarchive(file, extractDir)
}
