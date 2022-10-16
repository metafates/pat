package util

import (
	"github.com/metafates/pat/filesystem"
	"github.com/samber/lo"
	"os"
	"path/filepath"
	"strings"
)

func ResolveTilde(path string) string {
	path = filepath.Clean(path)
	// we do not care about cases where home dir is not present,
	// since pat is supposed to work with non-virtual environments,
	// so it's assumed that home dir is always present
	home := lo.Must(os.UserHomeDir())

	if path == "~" {
		return home
	} else if strings.HasPrefix(path, "~/") {
		return home + path[1:]
	}

	return path
}

func SubDirs(dir string) ([]string, error) {
	entries, err := filesystem.Api().ReadDir(dir)
	if err != nil {
		return nil, nil
	}

	return lo.FilterMap(entries, func(e os.FileInfo, _ int) (string, bool) {
		// check if symlink is a directory
		if e.Mode()&os.ModeSymlink != 0 {
			realPath, err := filepath.EvalSymlinks(filepath.Join(dir, e.Name()))
			if err != nil {
				return "", false
			}

			isDir, err := filesystem.Api().IsDir(realPath)
			if err != nil {
				return "", false
			}

			if isDir {
				return e.Name(), true
			}
		}

		if !e.IsDir() {
			return "", false
		}

		return e.Name(), true
	}), nil
}
