package tui

type state int

const (
	stateShellSelect state = iota
	statePathSelect
	statePathAdd
	stateEntriesPreview
	stateConfirmActions
	stateError
)
