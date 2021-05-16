package common

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
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

	if err := unpack(file, extractDir); err != nil {
		return err
	}

	return nil
}

func unpack(file, extractDir string) error {
	unfile := strings.TrimSuffix(file, ".gz")
	if err := ungzip(file, unfile); err != nil {
		return err
	}

	return untar(unfile, extractDir)
}

func untar(tarball, target string) error {
	reader, err := os.Open(tarball)
	if err != nil {
		return err
	}
	defer reader.Close()
	tarReader := tar.NewReader(reader)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		p := filepath.Join(target, header.Name)
		info := header.FileInfo()

		if info.IsDir() {
			if err = os.MkdirAll(p, info.Mode()); err != nil {
				return err
			}

			continue
		}

		file, err := os.OpenFile(p, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode())
		if err != nil {
			return err
		}
		defer file.Close()

		if _, err = io.Copy(file, tarReader); err != nil {
			return err
		}
	}

	return nil
}

func ungzip(source, target string) error {
	reader, err := os.Open(source)
	if err != nil {
		return err
	}

	defer reader.Close()

	archive, err := gzip.NewReader(reader)
	if err != nil {
		return err
	}

	defer archive.Close()

	target = filepath.Join(target, archive.Name)

	writer, err := os.Create(target)
	if err != nil {
		return err
	}

	defer writer.Close()

	_, err = io.Copy(writer, archive)

	return err
}
