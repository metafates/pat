package cmd

import (
	"github.com/metafates/pat/cli"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(cliCmd)
	cliCmd.PersistentFlags().StringVarP(&options.Shell, "shell", "s", "", "The shell to add the path to")
}

var cliCmd = &cobra.Command{
	Use:   "cli",
	Short: "Use pat as a cli tool",
}

var options = &cli.Options{}

func init() {
	cliCmd.AddCommand(cliAddCmd)
	cliAddCmd.Flags().StringVarP(&options.Path, "path", "p", "", "The path to add")
}

var cliAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a path to the shell",
	Run: func(cmd *cobra.Command, args []string) {
		options.Action = cli.ActionAdd
		handleErr(cli.Run(options))
	},
}

func init() {
	cliCmd.AddCommand(cliRemoveCmd)
	cliRemoveCmd.Flags().StringVarP(&options.Path, "path", "p", "", "The path to remove")
}

var cliRemoveCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove a path from the shell",
	Run: func(cmd *cobra.Command, args []string) {
		options.Action = cli.ActionRemove
		handleErr(cli.Run(options))
	},
}

func init() {
	cliCmd.AddCommand(cliListCmd)
}

var cliListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all paths in the shell",
	Run: func(cmd *cobra.Command, args []string) {
		options.Action = cli.ActionList
		handleErr(cli.Run(options))
	},
}

func init() {
	cliCmd.AddCommand(cliContainsCmd)
	cliContainsCmd.Flags().StringVarP(&options.Path, "path", "p", "", "The path to check")
}

var cliContainsCmd = &cobra.Command{
	Use:   "contains",
	Short: "Check if a path is in the shell",
	Run: func(cmd *cobra.Command, args []string) {
		options.Action = cli.ActionContains
		handleErr(cli.Run(options))
	},
}
