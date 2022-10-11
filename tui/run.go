package tui

import tea "github.com/charmbracelet/bubbletea"

func Run() error {
	return tea.NewProgram(NewModel(), tea.WithAltScreen()).Start()
}
