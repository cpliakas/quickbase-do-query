package cmd

import (
	"errors"

	"github.com/cpliakas/quickbase-do-query/cliutil"
	"github.com/cpliakas/quickbase-do-query/qb"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var varSetCfg *viper.Viper

var varSetCmd = &cobra.Command{
	Use:   "set [NAME] [VALUE]",
	Short: "Sets a database variable",
	Long:  ``,
	Args:  varSetCmdValidate,
	Run: func(cmd *cobra.Command, args []string) {

		input := &qb.SetVariableInput{
			AppID: globalCfg.AppID(),
			Name:  args[0],
			Value: args[1],
		}

		client := qb.NewClient(globalCfg)
		output, err := client.SetVariable(input)
		cliutil.HandleError(err, "error executing request")

		cliutil.PrintJSON(VarSetOutput{
			UserData: output.UserData,
			Name:     args[0],
			Value:    args[1],
		})
	},
}

func init() {
	varCmd.AddCommand(varSetCmd)
	varSetCfg = cliutil.InitConfig(qb.EnvVarPrefix)
}

func varSetCmdValidate(cmd *cobra.Command, args []string) error {
	globalCfg.RequireAppID = true
	if err := globalCfg.Validate(); err != nil {
		return err
	}

	if len(args) < 1 {
		return errors.New("missing required argument: [NAME]")
	}
	if len(args) < 2 {
		return errors.New("missing required argument: [VALUE]")
	}

	return nil
}

// VarSetOutput renders records in JSON.
type VarSetOutput struct {
	UserData string `json:"user_data,omitempty"`
	Name     string `json:"name"`
	Value    string `json:"value"`
}
