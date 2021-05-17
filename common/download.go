package common

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/go-enry/go-enry/v2"
)

type detector struct {
	binaries    []string
	libraries   []string
	manPages    []string
	completions map[string][]string
}

func (a *Asset) DownloadTo(destination *Destination) error {
	if a.Exists(destination) {
		return fmt.Errorf("%w: %s/%s", ErrAlreadyDownloaded, a.PackageName, a.Version)
	}

	root := destination.GetTmpDir()
	if err := ensureDir(root); err != nil {
		return err
	}

	root, err := ioutil.TempDir(root, "gr-")
	if err != nil {
		return err
	}

	defer os.RemoveAll(root)

	extractDir := destination.GetReleaseDirFor(a)
	if err := ensureDir(extractDir); err != nil {
		return err
	}

	file := path.Join(root, a.Name)
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

func (d *detector) detect(dir string) error {
	return filepath.Walk(dir, d.detectFile)
}

func (d *detector) detectFile(file string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	if info.IsDir() {
		return nil
	}

	data, err := readNFile(file, 1*1024*1024)
	if err != nil {
		return err
	}

	elfInfo, err := readElf(data)
	if err != nil {
		return err
	}

	lang := enry.GetLanguage(file, data)

	fmt.Println(file, isExec(info.Mode().Perm()), lang)
	fmt.Printf("%#v\n", elfInfo)

	return nil
}

func (d *detector) copyTo(destination *Destination) error {
	return nil
}
