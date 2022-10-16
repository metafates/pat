package cli

type Action int

const (
	ActionNone Action = iota
	ActionList
	ActionAdd
	ActionRemove
	ActionContains
)
