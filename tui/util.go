package tui

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/metafates/pat/color"
	"github.com/metafates/pat/filesystem"
	"github.com/metafates/pat/icon"
	"github.com/metafates/pat/path"
	"github.com/metafates/pat/util"
	"github.com/samber/mo"
	"path/filepath"
)

func (m *Model) pathInfo(p string) (info string) {
	m.toComplete = mo.None[string]()

	warning := func(s string) string {
		return lipgloss.NewStyle().Foreground(color.Yellow).Render(fmt.Sprintf("%s %s", icon.Warn, s))
	}

	invalid := func(s string) string {
		return lipgloss.NewStyle().Foreground(color.Red).Render(fmt.Sprintf("%s %s", icon.Cross, s))
	}

	valid := func(s string) string {
		return lipgloss.NewStyle().Foreground(color.Green).Render(fmt.Sprintf("%s %s", icon.Check, s))
	}

	suggestion := func(s string) string {
		return lipgloss.NewStyle().Foreground(color.Blue).Render(fmt.Sprintf("%s %s", icon.Info, s))
	}

	secondary := func(s string) string {
		return lipgloss.NewStyle().Faint(true).Render(s)
	}

	if !filepath.IsAbs(p) {
		asAbs := path.New(p).String()
		info = warning("Path is not absolute")
		info += "\n\n"
		info += secondary(asAbs)
		info += "\n\n"
		info += secondary(fmt.Sprintf("Press %s to complete", m.keymap.AcceptCompletion.Help().Key))
		m.toComplete = mo.Some[string](asAbs)
		return
	}

	exists, _ := filesystem.Api().Exists(p)

	if exists {
		isDir, _ := filesystem.Api().IsDir(p)
		if isDir {
			p := path.New(p)
			info = valid(fmt.Sprintf("Path exists, %s", util.Quantify(len(p.Executables()), "executable", "executables")))
		} else {
			info = invalid("Path exists but it's not a directory")
		}
	} else {
		dir := filepath.Dir(p)
		subDirs, err := util.SubDirs(dir)

		if err != nil {
			info = invalid("Path does not exist")
		} else {
			closest := util.FindClosest(filepath.Base(p), subDirs)
			if closest.IsPresent() {
				completion := filepath.Join(dir, closest.MustGet())
				info = suggestion(fmt.Sprintf("Did you mean %s?", completion))
				info += "\n\n"
				info += secondary(fmt.Sprintf("Press %s to complete", m.keymap.AcceptCompletion.Help().Key))
				m.toComplete = mo.Some(completion)
			} else {
				info = invalid("Path does not exist")
			}
		}
	}

	return
}
