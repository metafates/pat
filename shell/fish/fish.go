package fish

import (
	"fmt"
	"github.com/metafates/pat/constant"
	"github.com/metafates/pat/filesystem"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

var fishArrayElementRegex = regexp.MustCompile(`(?m).*?\|(.*?)\|$`)

type Fish struct{}

func New() *Fish {
	return &Fish{}
}

func (f *Fish) Name() string {
	return "fish"
}

func (f *Fish) cmd(code string) *exec.Cmd {
	return exec.Command(f.Name(), "-c", code)
}

func (f *Fish) writeFile(content string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	filename := fmt.Sprintf(".%s.%s", constant.App, f.Name())
	file, err := filesystem.Api().OpenFile(filepath.Join(home, filename), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return err
	}

	defer file.Close()

	content = fmt.Sprintf(`set -g PATH

%s

set -x PATH $PATH`, content)

	_, err = file.WriteString(content)
	return err
}

func (f *Fish) makeExport(path string) string {
	return fmt.Sprintf(`set -a PATH "%s"
`, path)
}

func (f *Fish) AddPath(path string) error {
	builder := strings.Builder{}

	paths, err := f.Paths()
	if err != nil {
		return err
	}

	builder.WriteString(f.makeExport(path))
	for _, p := range paths {
		builder.WriteString(f.makeExport(p))
	}

	return f.writeFile(builder.String())
}

func (f *Fish) RemovePath(path string) error {
	builder := strings.Builder{}

	paths, err := f.Paths()
	if err != nil {
		return err
	}

	for _, p := range paths {
		if p == path {
			continue
		}

		builder.WriteString(f.makeExport(p))
	}

	return f.writeFile(builder.String())
}

func (f *Fish) Paths() ([]string, error) {
	cmd := f.cmd("set -S PATH")
	out, err := cmd.Output()

	if err != nil {
		return nil, err
	}

	matched := fishArrayElementRegex.FindAllStringSubmatch(string(out), -1)
	paths := make([]string, len(matched))

	for i, m := range matched {
		paths[i] = m[1]
	}

	return paths, nil
}

func (f *Fish) Overwrite(paths []string) error {
	builder := strings.Builder{}

	for _, p := range paths {
		builder.WriteString(f.makeExport(p))
	}

	return f.writeFile(builder.String())
}

func (f *Fish) Available() bool {
	if _, err := exec.LookPath(f.Name()); err != nil {
		return false
	}

	return true
}
