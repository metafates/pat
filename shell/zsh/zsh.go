package zsh

import (
	"fmt"
	"github.com/metafates/pat/constant"
	"github.com/metafates/pat/where"
	"github.com/samber/lo"
	"os/exec"
	"strings"
)

type Zsh struct{}

func (z *Zsh) WhereScript() string {
	return where.ZshScript()
}

func (z *Zsh) Bin() string {
	return constant.Zsh
}

func (z *Zsh) CommentToken() string {
	return "#"
}

func (z *Zsh) Name() string {
	return "Zsh"
}

func New() *Zsh {
	return &Zsh{}
}

func (z *Zsh) cmd(code string) *exec.Cmd {
	return exec.Command(z.Name(), "-c", code)
}

func (z *Zsh) generateScript(content string) string {
	return fmt.Sprintf(`unset PATH

%s

export PATH`, content)
}

func (z *Zsh) makeExport(path string) string {
	return fmt.Sprintf(`PATH="%[1]s:$PATH"
`, path)
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

func (z *Zsh) GenerateExport(paths []string) string {
	builder := strings.Builder{}

	for i, _ := range paths {
		builder.WriteString(z.makeExport(paths[len(paths)-1-i]))
	}

	return z.generateScript(builder.String())
}

func (z *Zsh) Available() bool {
	if _, err := exec.LookPath(z.Name()); err != nil {
		return false
	}

	return true
}
