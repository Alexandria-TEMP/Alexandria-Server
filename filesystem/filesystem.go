package filesystem

import (
	"os"
	"path/filepath"
	"strconv"
)

type Filesystem struct {
	rootPath string
}

var cwd, _ = os.Getwd()
var defaultRootPath = filepath.Clean(filepath.Join(cwd, "vfs"))

func InitFilesystem() Filesystem {
	fs := Filesystem{
		rootPath: defaultRootPath,
	}

	err := os.MkdirAll(fs.rootPath, os.ModePerm)
	if err != nil {
		panic(err)
	}

	return fs
}

func (fs Filesystem) GetRepositoryPath(versionID uint, postID uint) string {
	return filepath.Join(fs.rootPath, strconv.FormatUint(uint64(postID), 10), strconv.FormatUint(uint64(versionID), 10))
}
