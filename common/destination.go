package common

import (
	"path"
)

type Destination struct {
	Root           string
	PackagesDir    string
	BinDir         string
	LibDir         string
	ManPagesDir    string
	TmpDir         string
	CompletionDirs map[string]string
}

func NewDestination(root string) *Destination {
	root = expand(root)

	return &Destination{
		Root:        root,
		PackagesDir: path.Join(root, "packages"),
		BinDir:      path.Join(root, "bin"),
		LibDir:      path.Join(root, "lib"),
		ManPagesDir: path.Join(root, "man"),
		TmpDir:      path.Join(root, "tmp"),
		CompletionDirs: map[string]string{
			"bash": path.Join(root, "bash.completions.d"),
			"zsh":  path.Join(root, "zsh.completions.d"),
			"fish": path.Join(root, "fish.completions.d"),
		},
	}
}

func (d *Destination) GetPackagesDir() string {
	return expand(d.PackagesDir)
}

func (d *Destination) GetTmpDir() string {
	return expand(d.TmpDir)
}

func (d *Destination) GetBinDir() string {
	return expand(d.BinDir)
}

func (d *Destination) GetLibDir() string {
	return expand(d.LibDir)
}

func (d *Destination) GetManPagesDir() string {
	return expand(d.ManPagesDir)
}

func (d *Destination) GetCompletionDir(shell string) string {
	c, ok := d.CompletionDirs[shell]
	if !ok {
		return ""
	}

	return expand(c)
}

func (d *Destination) GetPackageDirFor(a *Asset) string {
	return path.Join(d.GetPackagesDir(), a.PackageName, a.Version)
}
