package tui

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/metafates/pat/color"
	"github.com/metafates/pat/filesystem"
	"github.com/metafates/pat/icon"
	"github.com/metafates/pat/path"
	"github.com/metafates/pat/util"
	"github.com/samber/lo"
	"golang.org/x/exp/slices"
	"os"
	"path/filepath"
	"strings"
)

func (m *Model) View() string {
	switch m.state {
	case stateShellSelect:
		return m.styles.ListPadding.Render(m.shellSelectC.View())
	case statePathSelect:
		return m.styles.ListPadding.Render(m.pathSelectC.View())
	case stateEntriesPreview:
		return m.styles.ListPadding.Render(m.entriesPreviewC.View())
	case statePathAdd:
		return m.viewPathAdd()
	case stateConfirmActions:
		return m.viewConfirmActions()
	case stateError:
		return m.viewError()
	default:
		return ""
	}
}

func (m *Model) viewConfirmActions() string {
	builder := strings.Builder{}
	lines := []string{m.styles.Title.Render("Confirm"), ""}

	actions := lo.MapToSlice(m.onSave, func(k string, v action) lo.Tuple2[string, action] {
		return lo.Tuple2[string, action]{k, v}
	})

	slices.SortFunc(actions, func(a, b lo.Tuple2[string, action]) bool {
		return a.A < b.A
	})

	for _, t := range actions {
		p, a := t.A, t.B
		if a == actionNone {
			continue
		}

		builder.WriteString(lipgloss.NewStyle().Underline(true).Render(p))
		builder.WriteString(" ")
		builder.WriteString(icon.Arrow)
		builder.WriteString(" ")

		switch a {
		case actionAdd:
			builder.WriteString(lipgloss.NewStyle().Foreground(color.Green).Render("Will be added"))
		case actionRemove:
			builder.WriteString(lipgloss.NewStyle().Foreground(color.Red).Render("Will be removed"))
		default:
			panic("unknown action")
		}

		lines = append(lines, builder.String())
		builder.Reset()
	}

	if m.order.IsPresent() {
		lines = append(lines, lipgloss.NewStyle().Foreground(color.Yellow).Render("Order will be changed"))
	}

	if len(lines) == 0 {
		lines = append(lines, "Nothing to do")
	}

	return m.renderLines(true, lines...)
}

func (m *Model) viewPathAdd() string {
	value := util.ResolveTilde(m.textInputC.Value())
	exists, _ := filesystem.Api().Exists(value)
	var info string

	if exists {
		isDir, _ := filesystem.Api().IsDir(value)
		if isDir {
			p := path.New(value)
			info = lipgloss.NewStyle().Foreground(color.Green).Render(fmt.Sprintf("Path exists, %s", util.Quantify(len(p.Entries()), "entry", "entries")))
		} else {
			info = lipgloss.NewStyle().Foreground(color.Red).Render("Path exists but it's not a directory")
		}
	} else {
		dir := filepath.Dir(value)
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

			closest := util.FindClosest(filepath.Base(value), dirs)
			if closest.IsPresent() {
				info = lipgloss.NewStyle().Foreground(color.Yellow).Render("Did you mean " + filepath.Join(dir, closest.MustGet()) + "?")
			} else {
				info = lipgloss.NewStyle().Foreground(color.Red).Render("Path does not exist")
			}
		}
	}

	return m.renderLines(true,
		m.styles.Title.Render("Add path"),
		"",
		m.textInputC.View(),
		"",
		info,
	)
}

func (m *Model) viewError() string {
	return m.renderLines(true,
		m.styles.Title.Background(color.Red).Render("Error"),
		"",
		m.err.Error(),
	)
}

func (m *Model) renderLines(addHelp bool, lines ...string) string {
	height := len(lines)
	text := strings.Join(lines, "\n")

	if addHelp {
		text += strings.Repeat("\n", m.ttyHeight-height) + m.helpC.View(m.keymap)
	}

	return m.styles.Padding.Render(text)
}
