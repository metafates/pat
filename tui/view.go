package tui

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"strings"
)

func (m *Model) View() string {
	switch m.state {
	case stateShellSelect:
		return m.shellSelectC.View()
	case statePathSelect:
		return m.pathSelectC.View()
	case statePathAdd:
		return fmt.Sprintf("%t\n\n", m.textInputC.Focused()) + m.textInputC.View()
	case stateEntriesPreview:
		return m.entriesPreviewC.View()
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

	for p, a := range m.onSave {
		if a == actionNone {
			continue
		}

		builder.WriteString(p)
		builder.WriteString(" - ")

		switch a {
		case actionAdd:
			builder.WriteString("to add")
		case actionRemove:
			builder.WriteString("to remove")
		default:
			panic("unknown action")
		}

		builder.WriteString("\n")
	}

	if m.order.IsPresent() {
		builder.WriteString("Reorder PATH\n")
	}

	if builder.Len() == 0 {
		builder.WriteString("Nothing to do")
	} else {
		help := fmt.Sprintf("%s %s", m.keymap.Confirm.Help().Key, m.keymap.Confirm.Help().Desc)
		builder.WriteString(lipgloss.NewStyle().Faint(true).Render(help))
	}

	return builder.String()
}

func (m *Model) viewError() string {
	return m.err.Error()
}
