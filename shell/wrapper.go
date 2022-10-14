package shell

import (
	"github.com/metafates/pat/constant"
	"github.com/metafates/pat/filesystem"
	"github.com/metafates/pat/log"
	"github.com/samber/lo"
	"github.com/samber/mo"
	"github.com/spf13/afero"
	"os"
	"os/exec"
	"path/filepath"
)

type Wrapper struct {
	shell     instance
	isDefault mo.Option[bool]
	binPath   mo.Option[string]
}

func New(shell instance) *Wrapper {
	return &Wrapper{
		shell: shell,
	}
}
func (w *Wrapper) Available() bool {
	if _, err := exec.LookPath(w.shell.Bin()); err != nil {
		return false
	}

	return true
}

func (w *Wrapper) save(script string) error {
	w.backup()

	file, err := filesystem.Api().OpenFile(w.shell.WhereScript(), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return err
	}

	defer func(file afero.File) {
		err := file.Close()
		if err != nil {
			log.Warn(err)
		}
	}(file)

	err = constant.ScriptTemplate.Execute(file, struct {
		Comment string
		App     string
		Version string
		Script  string
	}{
		Comment: w.shell.CommentToken(),
		App:     constant.App,
		Version: constant.Version,
		Script:  script,
	})

	if err != nil {
		log.Error(err)
	}

	return err
}

func (w *Wrapper) AddPath(path string) error {
	paths, err := w.Paths()
	if err != nil {
		return err
	}

	// prepend to paths
	paths = append([]string{path}, paths...)

	return w.save(w.shell.GenerateExport(paths))
}

func (w *Wrapper) RemovePath(path string) error {
	paths, err := w.Paths()
	if err != nil {
		return err
	}

	paths = lo.Filter(paths, func(p string, _ int) bool {
		return p != path
	})

	return w.save(w.shell.GenerateExport(paths))
}

func (w *Wrapper) Export(paths []string) error {
	return w.save(w.shell.GenerateExport(paths))
}

func (w *Wrapper) Paths() ([]string, error) {
	return w.shell.Paths()
}

func (w *Wrapper) Name() string {
	return w.shell.Name()
}

func (w *Wrapper) IsDefault() bool {
	if w.isDefault.IsPresent() {
		return w.isDefault.MustGet()
	}

	defaultShell := filepath.Clean(os.Getenv("SHELL"))
	p, _ := w.BinPath()
	isDefault := p == defaultShell
	w.isDefault = mo.Some(isDefault)
	return isDefault
}

func (w *Wrapper) BinPath() (string, error) {
	if w.binPath.IsPresent() {
		return w.binPath.MustGet(), nil
	}

	p, err := exec.LookPath(w.shell.Bin())
	if err != nil {
		return "", err
	}

	w.binPath = mo.Some(p)
	return p, nil
}
