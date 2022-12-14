package bash

import (
	"fmt"
	"github.com/metafates/pat/constant"
	"github.com/metafates/pat/where"
	"github.com/samber/lo"
	"os/exec"
	"strings"
)

type Bash struct{}

func (b *Bash) WhereScript() string {
	return where.BashScript()
}

func (b *Bash) Bin() string {
	return constant.Bash
}

func (b *Bash) CommentToken() string {
	return "#"
}

func (b *Bash) Name() string {
	return "Bash"
}

func New() *Bash {
	return &Bash{}
}

func (b *Bash) cmd(code string) *exec.Cmd {
	return exec.Command(b.Name(), "-lc", code)
}

func (b *Bash) generateScript(content string) string {
	return fmt.Sprintf(`unset PATH

%s

export PATH`, content)
}

func (b *Bash) makeExport(path string) string {
	return fmt.Sprintf(`PATH="%s${PATH:+:${PATH}}"
`, path)
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

func (b *Bash) GenerateExport(paths []string) string {
	builder := strings.Builder{}

	for i, _ := range paths {
		builder.WriteString(b.makeExport(paths[len(paths)-1-i]))
	}

	return b.generateScript(builder.String())
}

func (b *Bash) Available() bool {
	if _, err := exec.LookPath(b.Name()); err != nil {
		return false
	}

	return true
}
