package cmd

import (
	"fmt"

	"github.com/cpliakas/quickbase-do-query/quickbase"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var varSetCmd = &cobra.Command{
	Use:   "set [NAME] [VALUE]",
	Short: "Sets a database variable",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			panic("required variables missing.")
		}

		cfg := quickbase.NewConfig()
		client := quickbase.NewClient(cfg)

		input := &quickbase.SetVariableInput{
			AppID: viper.GetString("app-id"),
			Name:  args[0],
			Value: args[1],
		}

		_, err := client.SetVariable(input)
		if err != nil {
			panic(err)
		}

		fmt.Printf("variable '%s' set: '%s'\n", args[0], args[1])
	},
}

func init() {
	varCmd.AddCommand(varSetCmd)
}
