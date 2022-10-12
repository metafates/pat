package tui

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/metafates/pat/path"
	"github.com/metafates/pat/shell"
	"github.com/samber/lo"
)

func (m *Model) Init() (cmd tea.Cmd) {
	if len(m.shells) == 1 {
		m.pushState(statePathSelect)
		cmd = m.pathSelectC.SetItems(
			lo.Map(lo.Must(m.shells[0].Paths()), func(p string, _ int) list.Item {
				return m.newItem(path.New(p))
			}),
		)
	} else {
		m.pushState(stateShellSelect)
		cmd = m.shellSelectC.SetItems(
			lo.Map(m.shells, func(s shell.Shell, _ int) list.Item {
				return m.newItem(s)
			}),
		)
	}

	return tea.Batch(cmd, m.textInputC.Focus(), textinput.Blink)
}
