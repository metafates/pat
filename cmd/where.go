package cmd

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/metafates/pat/color"
	"github.com/metafates/pat/constant"
	"github.com/metafates/pat/util"
	"github.com/metafates/pat/where"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

var wherePaths = []lo.Tuple2[string, func() string]{
	{"backup", where.Backup},
	{"config", where.Config},
	{"logs", where.Logs},
}

func init() {
	for _, p := range wherePaths {
		p.A = strings.ToLower(p.A)
	}
}

func init() {
	rootCmd.AddCommand(whereCmd)

	for _, n := range wherePaths {
		whereCmd.Flags().BoolP(n.A, string(n.A[0]), false, n.A+" path")
	}

	whereCmd.MarkFlagsMutuallyExclusive(lo.Map(wherePaths, func(t lo.Tuple2[string, func() string], _ int) string {
		return t.A
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
			if lo.Must(cmd.Flags().GetBool(n.A)) {
				cmd.Println(n.B())
				return
			}
		}

		for i, n := range wherePaths {
			cmd.Printf("%s %s\n", headerStyle(util.Capitalize(n.A)+"?"), yellowStyle("--"+n.A))
			cmd.Println(n.B())

			if i < len(wherePaths)-1 {
				cmd.Println()
			}
		}
	},
}
