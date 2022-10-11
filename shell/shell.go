package shell

import (
	"github.com/metafates/pat/log"
	"github.com/metafates/pat/shell/bash"
	"github.com/metafates/pat/shell/fish"
	"github.com/metafates/pat/shell/sh"
	"github.com/metafates/pat/shell/zsh"
	"github.com/samber/lo"
)

type Shell interface {
	Name() string
	AddPath(path string) error
	RemovePath(path string) error
	Paths() ([]string, error)
	Available() bool
}

func AvailableShells() []Shell {
	shells := lo.Filter([]Shell{
		fish.New(),
		zsh.New(),
		bash.New(),
		sh.New(),
	}, func(shell Shell, _ int) bool {
		return shell.Available()
	})

	log.Infof("found %d available shells", len(shells))
	return shells
}
