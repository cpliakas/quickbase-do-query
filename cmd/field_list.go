package cmd

import (
	"github.com/cpliakas/quickbase-do-query/cliutil"
	qb "github.com/cpliakas/quickbase-do-query/quickbase"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var fieldListCfg *viper.Viper

var fieldListCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists fields",
	Long:  ``,
	Args:  fieldListCmdValidate,
	Run: func(cmd *cobra.Command, args []string) {
		input := &qb.GetSchemaInput{ID: globalCfg.TableID()}

		client := qb.NewClient(globalCfg)
		output, err := client.GetSchema(input)
		cliutil.HandleError(err, "error executing request")

		// Build map of field ID to labels.
		fields := make(map[int]string)
		for _, f := range output.Fields {
			fields[f.FieldID] = f.Label
		}

		v := FieldListOutput{Fields: fields}
		cliutil.PrintJSON(v)
	},
}

func init() {
	fieldCmd.AddCommand(fieldListCmd)
	fieldListCfg = cliutil.InitConfig(qb.EnvVarPrefix)
}

func fieldListCmdValidate(cmd *cobra.Command, args []string) error {
	globalCfg.RequireTableID = true
	return globalCfg.Validate()
}

// FieldListOutput renders records in JSON.
type FieldListOutput struct {
	Fields map[int]string `json:"fields"`
}
