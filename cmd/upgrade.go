package cmd

import (
	"fmt"
	"github.com/metafates/pat/constant"
	"github.com/metafates/pat/upgrader"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(upgradeCmd)
}

var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: fmt.Sprintf("Upgrade %s to the latest version", constant.App),
	Run: func(cmd *cobra.Command, args []string) {
		handleErr(upgrader.Upgrade())
	},
}
