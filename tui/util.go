package tui

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/metafates/pat/color"
	"github.com/metafates/pat/filesystem"
	"github.com/metafates/pat/path"
	"github.com/metafates/pat/util"
	"github.com/samber/lo"
	"github.com/samber/mo"
	"os"
	"path/filepath"
)

func (m *Model) pathInfo(p string) (info string) {
	m.toComplete = mo.None[string]()

	isAbs := filepath.IsAbs(p)
	if !isAbs {
		info = lipgloss.NewStyle().Foreground(color.Yellow).Render("Path is not absolute")
		return
	}

	exists, _ := filesystem.Api().Exists(p)

	if exists {
		isDir, _ := filesystem.Api().IsDir(p)
		if isDir {
			p := path.New(p)
			info = lipgloss.NewStyle().Foreground(color.Green).Render(fmt.Sprintf("Path exists, %s", util.Quantify(len(p.Executables()), "executable", "executables")))
		} else {
			info = lipgloss.NewStyle().Foreground(color.Red).Render("Path exists but it's not a directory")
		}
	} else {
		dir := filepath.Dir(p)
		entries, err := filesystem.Api().ReadDir(dir)

		if err != nil {
			info = lipgloss.NewStyle().Foreground(color.Red).Render("Path does not exist")
		} else {
			dirs := lo.FilterMap(entries, func(e os.FileInfo, _ int) (string, bool) {
				// check if symlink is a directory
				if e.Mode()&os.ModeSymlink != 0 {
					realPath, err := filepath.EvalSymlinks(filepath.Join(dir, e.Name()))
					if err != nil {
						return "", false
					}

					isDir, err := filesystem.Api().IsDir(realPath)
					if err != nil {
						return "", false
					}

					if isDir {
						return e.Name(), true
					}
				}

				if !e.IsDir() {
					return "", false
				}

				return e.Name(), true
			})

			closest := util.FindClosest(filepath.Base(p), dirs)
			if closest.IsPresent() {
				completion := filepath.Join(dir, closest.MustGet())
				info = lipgloss.
					NewStyle().
					Foreground(color.Yellow).
					Render(fmt.Sprintf("Did you mean %s?", completion))
				info += lipgloss.NewStyle().Faint(true).Render(fmt.Sprintf(" Press %s to complete", m.keymap.AcceptCompletion.Help().Key))
				m.toComplete = mo.Some(completion)
			} else {
				info = lipgloss.NewStyle().Foreground(color.Red).Render("Path does not exist")
			}
		}
	}

	return
}
