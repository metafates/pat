package tui

import (
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/charmbracelet/lipgloss"
	"github.com/metafates/pat/color"
	"github.com/metafates/pat/icon"
	"github.com/metafates/pat/path"
	"github.com/metafates/pat/shell"
	"github.com/metafates/pat/util"
	"github.com/samber/lo"
	"path/filepath"
	"strings"
)

type item struct {
	model    *Model
	internal any
}

func (m *Model) newItem(internal any) *item {
	return &item{internal: internal, model: m}
}

func (i *item) marked() (markIcon string, marked bool) {
	switch x := i.internal.(type) {
	case *shell.Wrapper:
		marked = x.IsDefault()
		markIcon = lipgloss.NewStyle().Faint(true).Render("Default")
		return
	case *path.Path:
		pathAction, ok := i.model.getAction(x)
		if !ok {
			return
		}

		switch pathAction {
		case actionAdd:
			return lipgloss.NewStyle().Foreground(color.Green).Render(icon.Check), true
		case actionRemove:
			return lipgloss.NewStyle().Foreground(color.Red).Render(icon.Cross), true
		default:
			return
		}
	default:
		return
	}
}

func (i *item) FilterValue() string {
	switch i := i.internal.(type) {
	case *shell.Wrapper:
		return i.Name()
	case *path.Path:
		if i.IsDir() || !i.Exists() {
			return i.String()
		}

		return filepath.Base(i.String())
	default:
		panic("unknown type")
	}
}

func (i *item) Title() string {
	title := strings.Builder{}
	title.WriteString(i.FilterValue())

	if markIcon, marked := i.marked(); marked {
		title.WriteString(" ")
		switch i.internal.(type) {
		case *shell.Wrapper:
			title.WriteString(markIcon)
		case *path.Path:
			title.WriteString(markIcon)
		}
	}

	return title.String()
}

func (i *item) Description() string {
	switch i := i.internal.(type) {
	case *shell.Wrapper:
		return lo.Must(i.BinPath())
	case *path.Path:
		if !i.Exists() {
			return "Nonexistent"
		}

		if i.IsDir() {
			entries := util.Quantify(len(i.Entries()), "entry", "entries")
			size := i.SizeHuman()

			return fmt.Sprintf("%s, %s", entries, size)
		}

		if !i.IsExecutable() {
			return fmt.Sprintf("%s, Nonexecutable", i.SizeHuman())
		}

		return i.SizeHuman()
	default:
		panic("unknown type")
	}
}

func (i *item) Copy() {
	_ = clipboard.WriteAll(i.FilterValue())
}
