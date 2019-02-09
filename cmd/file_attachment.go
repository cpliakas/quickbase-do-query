package cmd

import (
	"github.com/spf13/cobra"
)

var fileAttachmentCmd = &cobra.Command{
	Use:   "file-attachment",
	Short: "Commands that act on file attachment",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(fileAttachmentCmd)
}
