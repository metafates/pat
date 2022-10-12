package fish

import (
	"fmt"
	"github.com/metafates/pat/log"
	"github.com/samber/lo"
	"os/exec"
	"regexp"
	"strings"
)

var fishArrayElementRegex = regexp.MustCompile(`(?m).*?\|(.*?)\|$`)

type Fish struct{}

func New() *Fish {
	return &Fish{}
}

func (f *Fish) cmd(code string) *exec.Cmd {
	cmd := exec.Command(f.Name(), "--command", code)
	log.Infof("executing %s", code)
	return cmd
}

func (f *Fish) Name() string {
	return "fish"
}

func (f *Fish) AddPath(path string) error {
	return f.
		cmd(
			fmt.Sprintf(`
set --universal fish_user_paths "%s" $fish_user_paths 
`, path),
		).
		Run()
}

func (f *Fish) RemovePath(path string) error {
	return f.
		cmd(
			fmt.Sprintf(`
if set -l index (contains -i "%s" $fish_user_paths)
	set --erase --universal fish_user_paths[$index]
end
`, path),
		).
		Run()
}

func (f *Fish) Overwrite(paths []string) error {
	builder := strings.Builder{}

	for _, p := range paths {
		builder.WriteString(fmt.Sprintf(`"%s"`, p))
		builder.WriteString("")
	}

	return f.cmd("set --universal fish_user_paths " + builder.String()).Run()
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
