package util

import (
	"github.com/samber/lo"
	"os"
	"strings"
)

func ResolveTilde(path string) string {
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
