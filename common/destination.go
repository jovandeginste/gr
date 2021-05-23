package common

import (
	"fmt"
	"path"
)

type Destination struct {
	Root           string
	PackagesDir    string
	SourcesDir     string
	BinDir         string
	LibDir         string
	ManPagesDir    string
	TmpDir         string
	CompletionDirs map[string]string
}

func (d *Destination) GetAllManDirs() []string {
	var r []string

	for i := 1; i < 9; i++ {
		r = append(r, path.Join(d.ManPagesDir, fmt.Sprintf("man%d", i)))
	}

	return r
}

func NewLinuxRootDestination(root string) *Destination {
	return &Destination{
		Root:        root,
		PackagesDir: path.Join(root, "var", "lib", "gr"),
		BinDir:      path.Join(root, "usr", "bin"),
		LibDir:      path.Join(root, "usr", "lib"),
		ManPagesDir: path.Join(root, "usr", "share", "man"),
		TmpDir:      path.Join(root, "tmp"),
		CompletionDirs: map[string]string{
			"bash": path.Join(root, "usr", "share", "bash-completion", "completions"),
			"zsh":  path.Join(root, "usr", "share", "zsh", "vendor-completions"),
			"fish": path.Join(root, "usr", "share", "fish", "vendor_completions.d"),
		},
	}
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
			"bash": path.Join(root, "completions.d", "bash"),
			"zsh":  path.Join(root, "completions.d", "zsh"),
			"fish": path.Join(root, "completions.d", "fish"),
		},
	}
}

func (d *Destination) EnsureDirs() error {
	dirs := []string{
		d.GetPackagesDir(),
		d.GetTmpDir(),
		d.GetBinDir(),
		d.GetLibDir(),
		d.GetManPagesDir(),
	}

	dirs = append(dirs, d.GetAllManDirs()...)
	dirs = append(dirs, d.GetAllCompletionDirs()...)

	for _, d := range dirs {
		if err := ensureDir(d); err != nil {
			return err
		}
	}

	return nil
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

func (d *Destination) GetAllCompletionDirs() []string {
	res := []string{}

	for k := range d.CompletionDirs {
		res = append(res, d.GetCompletionDir(k))
	}

	return res
}

func (d *Destination) GetCompletionDir(shell string) string {
	c, ok := d.CompletionDirs[shell]
	if !ok {
		return ""
	}

	return expand(c)
}

func (d *Destination) GetPackageDirFor(name, version string) string {
	return path.Join(d.GetPackagesDir(), name, version)
}

func (d *Destination) GetReleaseDirFor(a *Asset) string {
	return path.Join(d.GetPackageDirFor(a.PackageName, a.Version), "release")
}

func (d *Destination) GetSourceDirFor(r *Release) string {
	return path.Join(d.GetPackageDirFor(r.PackageName, r.Version), "source")
}
