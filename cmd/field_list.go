package cmd

import (
	"github.com/cpliakas/quickbase-do-query/qbutil"

	"github.com/cpliakas/quickbase-do-query/cliutil"
	qb "github.com/cpliakas/quickbase-do-query/quickbase"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var fieldListCfg *viper.Viper

// TODO: Replace panics with a better error handling mechanism.
var fieldListCmd = &cobra.Command{
	Use:     "list",
	Short:   "Lists fields",
	Long:    ``,
	PreRunE: globalCfg.PreRunE,
	Run: func(cmd *cobra.Command, args []string) {
		qbutil.RequireTableID(globalCfg)

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

// FieldListOutput renders records in JSON.
type FieldListOutput struct {
	Fields map[int]string `json:"fields"`
}
