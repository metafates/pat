package fish

import (
	"fmt"
	"github.com/metafates/pat/constant"
	"github.com/metafates/pat/where"
	"os/exec"
	"regexp"
	"strings"
)

var fishArrayElementRegex = regexp.MustCompile(`(?m).*?\|(.*?)\|$`)

type Fish struct{}

func (f *Fish) WhereScript() string {
	return where.FishScript()
}

func (f *Fish) CommentToken() string {
	return "#"
}

func New() *Fish {
	return &Fish{}
}

func (f *Fish) Name() string {
	return "Fish"
}

func (f *Fish) Bin() string {
	return constant.Fish
}

func (f *Fish) cmd(code string) *exec.Cmd {
	return exec.Command(f.Name(), "-c", code)
}

func (f *Fish) generateScript(content string) string {
	return fmt.Sprintf(`set --global PATH

%s

set --export PATH $PATH`, content)
}

func (f *Fish) makeExport(path string) string {
	return fmt.Sprintf(`set -a PATH "%s"
`, path)
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

func (f *Fish) GenerateExport(paths []string) string {
	builder := strings.Builder{}

	for _, p := range paths {
		builder.WriteString(f.makeExport(p))
	}

	return f.generateScript(builder.String())
}
