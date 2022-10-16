package shell

import (
	"github.com/metafates/pat/log"
	"github.com/metafates/pat/shell/bash"
	"github.com/metafates/pat/shell/fish"
	"github.com/metafates/pat/shell/zsh"
	"github.com/samber/lo"
)

type instance interface {
	Name() string
	Bin() string
	WhereScript() string
	CommentToken() string
	GenerateExport(paths []string) (script string)
	Paths() ([]string, error)
}

func AvailableShells() []*Wrapper {
	shells := lo.Filter([]*Wrapper{
		New(fish.New()),
		New(zsh.New()),
		New(bash.New()),
	}, func(shell *Wrapper, _ int) bool {
		return shell.Available()
	})

	log.Infof("found %d available shells", len(shells))
	return shells
}

func Get(name string) (wrapper *Wrapper, ok bool) {
	wrapper, ok = lo.Find(AvailableShells(), func(w *Wrapper) bool {
		return w.Bin() == name
	})
	return
}
