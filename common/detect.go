package common

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/go-enry/go-enry/v2"
	"github.com/h2non/filetype/types"
)

var shells = []string{"bash", "fish", "zsh"}

type detector struct {
	binaries    []string
	libraries   []string
	manPages    []string
	completions map[string][]string
}

func (r *Release) Detect(destination *Destination) error {
	packageDir := destination.GetPackageDirFor(r.PackageName, r.Version)

	d := detector{
		completions: map[string][]string{},
	}
	if err := d.detect(packageDir); err != nil {
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

	if !info.Mode().IsRegular() {
		return nil
	}

	data, err := readNFile(file, 8192)
	if err != nil {
		return err
	}

	elfInfo, err := readElf(data)
	if err != nil {
		return err
	}

	lang := enry.GetLanguage(file, data)

	d.addToCatalog(file, lang, elfInfo)

	return nil
}

func (d *detector) addToCatalog(file string, lang string, elfInfo types.Type) {
	if elfInfo.MIME.Value == "application/x-executable" {
		d.binaries = append(d.binaries, file)

		return
	}

	if lang == "Roff Manpage" {
		d.manPages = append(d.manPages, file)

		return
	}

	if lang == "Shell" {
		if strings.Contains(file, "completion") {
			for _, s := range shells {
				if strings.Contains(file, s) {
					d.completions[s] = append(d.completions[s], file)
				}
			}
		}
	}

	fmt.Println(file, lang)
	fmt.Printf("%#v\n", elfInfo)
}

func (d *detector) copyTo(destination *Destination) error {
	fmt.Printf("%#v\n", d)
	if err := linkAll(d.binaries, destination.GetBinDir()); err != nil {
		return err
	}

	if err := linkAll(d.manPages, destination.GetManPagesDir()); err != nil {
		return err
	}

	if err := linkAll(d.libraries, destination.GetLibDir()); err != nil {
		return err
	}

	for k, v := range d.completions {
		if err := linkAll(v, destination.GetCompletionDir(k)); err != nil {
			return err
		}
	}

	return nil
}

func linkAll(files []string, destination string) error {
	for _, b := range files {
		base := path.Base(b)
		target := path.Join(destination, base)

		if err := os.Remove(target); err != nil && !os.IsNotExist(err) {
			return err
		}

		if err := os.Link(b, target); err != nil {
			return err
		}
	}

	return nil
}
