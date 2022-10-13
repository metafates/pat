package zsh

import (
	"fmt"
	"github.com/metafates/pat/constant"
	"github.com/metafates/pat/filesystem"
	"github.com/samber/lo"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Zsh struct{}

func (z *Zsh) Name() string {
	return "zsh"
}

func New() *Zsh {
	return &Zsh{}
}

func (z *Zsh) cmd(code string) *exec.Cmd {
	return exec.Command(z.Name(), "-c", code)
}

func (z *Zsh) writeFile(content string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	filename := fmt.Sprintf(".%s.%s", constant.App, z.Name())
	file, err := filesystem.Api().OpenFile(filepath.Join(home, filename), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return err
	}

	defer file.Close()

	content = fmt.Sprintf(`unset PATH

%s

export PATH`, content)

	_, err = file.WriteString(content)
	return err
}

func (z *Zsh) makeExport(path string) string {
	return fmt.Sprintf(`PATH="%[1]s:$PATH"
`, path)
}

func (z *Zsh) AddPath(path string) error {
	builder := strings.Builder{}

	paths, err := z.Paths()
	if err != nil {
		return err
	}

	for i, _ := range paths {
		builder.WriteString(z.makeExport(paths[len(paths)-1-i]))
	}

	builder.WriteString(z.makeExport(path))

	return z.writeFile(builder.String())
}

func (z *Zsh) RemovePath(path string) error {
	builder := strings.Builder{}

	paths, err := z.Paths()
	if err != nil {
		return err
	}

	for i, _ := range paths {
		p := paths[len(paths)-1-i]
		if p == path {
			continue
		}

		builder.WriteString(z.makeExport(p))
	}

	return z.writeFile(builder.String())
}

func (z *Zsh) Paths() ([]string, error) {
	cmd := z.cmd("echo $PATH")
	out, err := cmd.Output()

	if err != nil {
		return nil, err
	}

	paths := lo.Filter(strings.Split(strings.TrimSpace(string(out)), ":"), func(s string, _ int) bool {
		return s != ""
	})

	return paths, nil
}

func (z *Zsh) Overwrite(paths []string) error {
	builder := strings.Builder{}

	for i, _ := range paths {
		builder.WriteString(z.makeExport(paths[len(paths)-1-i]))
	}

	return z.writeFile(builder.String())
}

func (z *Zsh) Available() bool {
	if _, err := exec.LookPath(z.Name()); err != nil {
		return false
	}

	return true
}
