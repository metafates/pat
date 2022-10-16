package cli

import (
	"errors"
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/metafates/pat/color"
	"github.com/metafates/pat/shell"
	"github.com/samber/lo"
	"os"
)

func Run(options *Options) error {
	wrapper, exists := shell.Get(options.Shell)
	if !exists {
		return fmt.Errorf("shell %s is not supported", lipgloss.NewStyle().Foreground(color.Yellow).Render(options.Shell))
	}

	switch options.Action {
	case ActionNone:
		return errors.New("no action specified")
	case ActionAdd:
		return wrapper.AddPath(options.Path)
	case ActionRemove:
		return wrapper.RemovePath(options.Path)
	case ActionContains:
		paths, err := wrapper.Paths()
		if err != nil {
			return err
		}

		if lo.Contains(paths, options.Path) {
			fmt.Println("true")
		} else {
			fmt.Println("false")
			os.Exit(1)
		}
	case ActionList:
		paths, err := wrapper.Paths()
		if err != nil {
			return err
		}

		for _, path := range paths {
			fmt.Println(path)
		}
	default:
		panic("unreachable")
	}

	return nil
}
