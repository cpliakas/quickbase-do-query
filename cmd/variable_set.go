package cmd

import (
	"fmt"

	"github.com/cpliakas/quickbase-do-query/cliutil"
	"github.com/cpliakas/quickbase-do-query/qbutil"
	qb "github.com/cpliakas/quickbase-do-query/quickbase"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var varSetCmd = &cobra.Command{
	Use:   "set [NAME] [VALUE]",
	Short: "Sets a database variable",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		globalCfg.InitConfig()
		qbutil.RequireAppID(globalCfg)

		input := &qb.SetVariableInput{
			AppID: viper.GetString("app-id"),
			Name:  args[0],
			Value: args[1],
		}

		client := qb.NewClient(globalCfg)
		_, err := client.SetVariable(input)
		cliutil.HandleError(err, "error executing request")

		fmt.Printf("variable '%s' set: '%s'\n", args[0], args[1])
	},
}

func init() {
	varCmd.AddCommand(varSetCmd)
}
