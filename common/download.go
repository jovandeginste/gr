package common

import (
	"io/ioutil"
	"os"
	"path"

	archiver "github.com/mholt/archiver/v3"
)

type Destination struct {
	Root        string
	Binaries    string
	Libraries   string
	ManPages    string
	Completions map[string]string
}

type detector struct {
	binaries    []string
	libraries   []string
	manPages    []string
	completions map[string][]string
}

func (a *Asset) DownloadTo(destination *Destination) error {
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

	a.Logger.Infof("Unpacking in '%s'...", extractDir)
	if err := unpack(file, extractDir); err != nil {
		return err
	}

	d := detector{}
	if err := d.detect(extractDir); err != nil {
		return err
	}

	if err := d.copyTo(destination); err != nil {
		return err
	}

	return nil
}

func unpack(file, extractDir string) error {
	return archiver.Unarchive(file, extractDir)
}

func (d *detector) detect(dir string) error {
	return nil
}

func (d *detector) copyTo(destination *Destination) error {
	return nil
}
