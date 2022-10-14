package tui

type state int

const (
	stateShellSelect state = iota
	statePathSelect
	statePathAdd
	stateExecutablesPreview
	stateConfirmActions
	stateError
)
