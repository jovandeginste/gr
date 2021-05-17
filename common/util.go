package common

import (
	"os"

	"github.com/gabriel-vasile/mimetype"
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

	_, err = file.Read(byteSlice)

	return byteSlice, err
}

func readElf(file string) (*mimetype.MIME, error) {
	return mimetype.DetectFile(file)
}
