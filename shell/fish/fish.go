package fish

import (
	"fmt"
	"github.com/samber/lo"
	"os/exec"
	"regexp"
)

var fishArrayElementRegex = regexp.MustCompile(`(?m).*?\|(.*?)\|$`)

type Fish struct{}

func New() *Fish {
	return &Fish{}
}

func (f *Fish) cmd(code string) *exec.Cmd {
	return exec.Command(f.Name(), "--command", code)
}

func (f *Fish) Name() string {
	return "fish"
}

func (f *Fish) AddPath(path string) error {
	return f.
		cmd(fmt.Sprintf(`fish_add_path %s`, path)).
		Run()
}

func (f *Fish) RemovePath(path string) error {
	p, err := f.Paths()
	if err != nil {
		return err
	}

	return f.
		cmd(
			fmt.Sprintf(`
if set -l index (contains -i "%s" "%s")
	set --erase --universal fish_user_paths[$index]
end
`, path, p),
		).
		Run()
}

func (f *Fish) Paths() ([]string, error) {
	cmd := f.cmd("set -S fish_user_paths")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	paths := fishArrayElementRegex.FindAllStringSubmatch(string(out), -1)

	return lo.Map(paths, func(p []string, _ int) string {
		return p[1]
	}), nil
}

func (f *Fish) Available() bool {
	if _, err := exec.LookPath(f.Name()); err != nil {
		return false
	}

	return true
}
