package cmd

import (
	"fmt"

	"github.com/cpliakas/quickbase-do-query/cliutil"
	"github.com/cpliakas/quickbase-do-query/qbutil"
	qb "github.com/cpliakas/quickbase-do-query/quickbase"
	"github.com/spf13/cobra"
)

var varSetCmd = &cobra.Command{
	Use:   "set [NAME] [VALUE]",
	Short: "Sets a database variable",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		globalCfg.InitConfig()
		qbutil.RequireAppID(globalCfg)

		cliutil.RequireArg(args, 0, "name")
		cliutil.RequireArg(args, 1, "value")

		input := &qb.SetVariableInput{
			AppID: globalCfg.AppID(),
			Name:  args[0],
			Value: args[1],
		}

		client := qb.NewClient(globalCfg)
		_, err := client.SetVariable(input)
		cliutil.HandleError(err, "error executing request")

		// TODO return JSON.
		fmt.Printf("variable '%s' set: '%s'\n", args[0], args[1])
	},
}

func init() {
	varCmd.AddCommand(varSetCmd)
}
