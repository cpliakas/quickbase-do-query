package cmd

import (
	"fmt"
	"os"

	"github.com/cpliakas/quickbase-do-query/qbutil"

	"github.com/cpliakas/quickbase-do-query/cliutil"
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
	flags := cliutil.NewFlagger(rootCmd, cfg)

	flags.PersistentString("app-id", "I", "", "application's dbid")
	flags.PersistentString("app-token", "A", "", "app token used with ticket to to authenticate API requests")
	flags.PersistentString("realm-host", "R", "", "The realm host, e.g., 'https://MYREALM.quickbase.com'")
	flags.PersistentString("ticket", "T", "", "ticket used to authenticate API requests")
	flags.PersistentString("table-id", "t", "", "table's dbid")
	flags.PersistentString("user-token", "U", "", "user token used to authenticate API requests")

	globalCfg = qbutil.NewGlobalConfig(cfg)
}
