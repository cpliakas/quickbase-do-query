package cmd

import (
	"github.com/spf13/cobra"
)

var varCmd = &cobra.Command{
	Use:   "var",
	Short: "Commands that act on database variables",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(varCmd)
}
