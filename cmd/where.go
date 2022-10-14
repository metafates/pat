package cmd

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/metafates/pat/color"
	"github.com/metafates/pat/constant"
	"github.com/metafates/pat/where"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

type whereTarget struct {
	Name     string
	ArgShort string
	ArgFull  string
	Where    func() string
}

var wherePaths = []whereTarget{
	{"Zsh Script", "", "zsh", where.ZshScript},
	{"Fish Script", "", "fish", where.FishScript},
	{"Bash Script", "", "bash", where.BashScript},

	{"Backup", "b", "backup", where.Backup},
	{"Config", "c", "config", where.Config},
	{"Logs", "l", "logs", where.Logs},
}

func init() {
	for _, p := range wherePaths {
		p.ArgFull = strings.ToLower(p.ArgFull)
		p.ArgShort = strings.ToLower(p.ArgShort)
	}
}

func init() {
	rootCmd.AddCommand(whereCmd)

	for _, n := range wherePaths {
		if n.ArgShort != "" {
			whereCmd.Flags().BoolP(n.ArgFull, n.ArgShort, false, n.Name+" path")
		} else {
			whereCmd.Flags().Bool(n.ArgFull, false, n.Name+" path")
		}
	}

	whereCmd.MarkFlagsMutuallyExclusive(lo.Map(wherePaths, func(t whereTarget, _ int) string {
		return t.ArgFull
	})...)

	whereCmd.SetOut(os.Stdout)
}

var whereCmd = &cobra.Command{
	Use:   "where",
	Short: "Show the paths for a files related to the " + constant.App,
	Run: func(cmd *cobra.Command, args []string) {
		headerStyle := lipgloss.NewStyle().Bold(true).Foreground(color.HiPurple).Render
		yellowStyle := lipgloss.NewStyle().Foreground(color.Yellow).Render

		for _, n := range wherePaths {
			if lo.Must(cmd.Flags().GetBool(n.ArgFull)) {
				cmd.Println(n.Where())
				return
			}
		}

		for i, n := range wherePaths {
			cmd.Printf("%s %s\n", headerStyle(n.Name+"?"), yellowStyle("--"+n.ArgFull))
			cmd.Println(n.Where())

			if i < len(wherePaths)-1 {
				cmd.Println()
			}
		}
	},
}
