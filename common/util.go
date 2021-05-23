package common

import (
	"os"
	"path"
	"strings"

	"github.com/h2non/filetype"
	"github.com/h2non/filetype/types"
	"github.com/mholt/archiver/v3"
	"github.com/mitchellh/go-homedir"
)

func expand(dir string) string {
	dir, _ = homedir.Expand(dir) //nolint:errcheck

	return dir
}

func isAppImage(file string) bool {
	return strings.HasSuffix(file, ".AppImage")
}

func isArchive(file string) bool {
	_, err := archiver.ByExtension(file)

	return err == nil
}

func moveAppImage(file, extractDir string) error {
	b := path.Base(file)
	b = strings.TrimSuffix(b, ".AppImage")

	target := path.Join(extractDir, b)

	if err := os.Chmod(file, 0o700); err != nil {
		return err
	}

	return os.Rename(file, target)
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

func readNFile(p string, n uint32) ([]byte, error) {
	file, err := os.Open(p)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	// Read up to len(b) bytes from the File
	// Zero bytes written means end of file
	// End of file returns error type io.EOF

	byteSlice := make([]byte, n)

	size, err := file.Read(byteSlice)

	return byteSlice[:size], err
}

func readElf(data []byte) (types.Type, error) {
	return filetype.Match(data)
}
