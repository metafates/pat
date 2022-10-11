package cmd

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	cc "github.com/ivanpirog/coloredcobra"
	"github.com/metafates/pat/color"
	"github.com/metafates/pat/constant"
	"github.com/metafates/pat/filesystem"
	"github.com/metafates/pat/icon"
	"github.com/metafates/pat/log"
	"github.com/metafates/pat/tui"
	"github.com/metafates/pat/where"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     strings.ToLower(constant.App),
	Short:   "App description",
	Version: constant.Version,
	Run: func(cmd *cobra.Command, args []string) {
		handleErr(tui.Run())
	},
}

func Execute() {
	cc.Init(&cc.Config{
		RootCmd:       rootCmd,
		Headings:      cc.HiCyan + cc.Bold + cc.Underline,
		Commands:      cc.HiYellow + cc.Bold,
		Example:       cc.Italic,
		ExecName:      cc.Bold,
		Flags:         cc.Bold,
		FlagsDataType: cc.Italic + cc.HiBlue,
	})

	// Clears temp files on each run.
	// It should not affect startup time since it's being run in parallel.
	go func() {
		_ = filesystem.Api().RemoveAll(where.Temp())
	}()

	_ = rootCmd.Execute()
}

func handleErr(err error) {
	if err != nil {
		log.Error(err)
		_, _ = fmt.Fprintf(
			os.Stderr,
			"%s %s\n",
			lipgloss.NewStyle().Foreground(color.Red).Render(icon.Cross),
			strings.Trim(err.Error(), " \n"),
		)
		os.Exit(1)
	}
}
