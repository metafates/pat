package tui

func (m *Model) View() string {
	switch m.state {
	case stateShellSelect:
		return m.shellSelectC.View()
	case statePathSelect:
		return m.pathSelectC.View()
	case statePathAdd:
		return m.textInputC.View()
	case stateError:
		return m.viewError()
	case stateQuit:
		return "Goodbye!"
	default:
		return ""
	}
}

func (m *Model) viewError() string {
	return m.err.Error()
}
