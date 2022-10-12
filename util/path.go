package util

import (
	"os"
	"strings"
)

func ResolveTilde(path string) (string, error) {
	dir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	if path == "~" {
		return dir, nil
	} else if strings.HasPrefix(path, "~/") {
		return dir + path[1:], nil
	}

	return path, nil
}
