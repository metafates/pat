package bash

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

type Bash struct{}

func (b *Bash) Name() string {
	return "bash"
}

func New() *Bash {
	return &Bash{}
}

func (b *Bash) cmd(code string) *exec.Cmd {
	return exec.Command(b.Name(), "-c", code)
}

func (b *Bash) writeFile(content string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	filename := fmt.Sprintf(".%s.%s", constant.App, b.Name())
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

func (b *Bash) makeExport(path string) string {
	return fmt.Sprintf(`PATH="%[1]s:$PATH"
`, path)
}

func (b *Bash) AddPath(path string) error {
	builder := strings.Builder{}

	paths, err := b.Paths()
	if err != nil {
		return err
	}

	for i, _ := range paths {
		builder.WriteString(b.makeExport(paths[len(paths)-1-i]))
	}

	builder.WriteString(b.makeExport(path))

	return b.writeFile(builder.String())
}

func (b *Bash) RemovePath(path string) error {
	builder := strings.Builder{}

	paths, err := b.Paths()
	if err != nil {
		return err
	}

	for i, _ := range paths {
		p := paths[len(paths)-1-i]
		if p == path {
			continue
		}

		builder.WriteString(b.makeExport(p))
	}

	return b.writeFile(builder.String())
}

func (b *Bash) Paths() ([]string, error) {
	cmd := b.cmd("echo $PATH")
	out, err := cmd.Output()

	if err != nil {
		return nil, err
	}

	paths := lo.Filter(strings.Split(strings.TrimSpace(string(out)), ":"), func(s string, _ int) bool {
		return s != ""
	})

	return paths, nil
}

func (b *Bash) Overwrite(paths []string) error {
	builder := strings.Builder{}

	for i, _ := range paths {
		builder.WriteString(b.makeExport(paths[len(paths)-1-i]))
	}

	return b.writeFile(builder.String())
}

func (b *Bash) Available() bool {
	if _, err := exec.LookPath(b.Name()); err != nil {
		return false
	}

	return true
}
