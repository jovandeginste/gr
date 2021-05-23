package common

import (
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/go-enry/go-enry/v2"
	"github.com/h2non/filetype/types"
	"github.com/sirupsen/logrus"
)

var shells = []string{"bash", "fish", "zsh"}

type detector struct {
	name        string
	logger      *logrus.Logger
	binaries    []string
	libraries   []string
	manPages    map[string][]string
	completions map[string][]string
}

func (r *Release) Detect(destination *Destination) error {
	packageDir := destination.GetPackageDirFor(r.PackageName, r.Version)

	d := detector{
		name:        r.PackageName,
		logger:      r.Logger,
		manPages:    map[string][]string{},
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
		d.addManPage(file)

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

	d.logger.Tracef("Language of '%s': '%s'", file, lang)
	d.logger.Tracef("Binary information for '%s': %#v", file, elfInfo)
}

func (d *detector) addManPage(file string) {
	ss := strings.Split(file, ".")
	r := "man1"

	for i := len(ss) - 1; i > 0; i-- {
		s := ss[i]
		if n, err := strconv.Atoi(s); err == nil {
			if n > 0 && n < 10 {
				r = "man" + s

				break
			}
		}
	}

	d.manPages[r] = append(d.manPages[r], file)
}

func (d *detector) copyTo(destination *Destination) error {
	if err := d.linkAll(d.binaries, destination.GetBinDir(), ""); err != nil {
		return err
	}

	for k, v := range d.manPages {
		if err := d.linkAll(v, path.Join(destination.GetManPagesDir(), k), ""); err != nil {
			return err
		}
	}

	if err := d.linkAll(d.libraries, destination.GetLibDir(), ""); err != nil {
		return err
	}

	for k, v := range d.completions {
		if err := d.linkAll(v, destination.GetCompletionDir(k), d.name+"."); err != nil {
			return err
		}
	}

	return nil
}

func (d *detector) linkAll(files []string, destination string, prefix string) error {
	for _, b := range files {
		base := prefix + path.Base(b)
		target := path.Join(destination, base)

		if err := os.Remove(target); err != nil && !os.IsNotExist(err) {
			return err
		}

		d.logger.Infof("Linking: '%s' to '%s'", b, target)
		if err := os.Link(b, target); err != nil {
			return err
		}
	}

	return nil
}
