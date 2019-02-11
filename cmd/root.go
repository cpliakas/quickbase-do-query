package cmd

import (
	"fmt"
	"os"

	"github.com/cpliakas/quickbase-do-query/cliutil"
	"github.com/cpliakas/quickbase-do-query/qbutil"
	qb "github.com/cpliakas/quickbase-do-query/quickbase"
	"github.com/spf13/cobra"
)

var globalCfg qbutil.GlobalConfig

var rootCmd = &cobra.Command{
	Use:   "quickbase-do-query",
	Short: "A command line tool that gets records from a Quick Base table.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cfg := cliutil.InitConfig(qb.EnvVarPrefix)
	globalCfg = qbutil.NewGlobalConfig(rootCmd, cfg)
}
