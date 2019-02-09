package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/cpliakas/quickbase-do-query/qbutil"

	"github.com/cpliakas/quickbase-do-query/cliutil"
	qb "github.com/cpliakas/quickbase-do-query/quickbase"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var fieldListCfg *viper.Viper

// TODO: Replace panics with a better error handling mechanism.
var fieldListCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists fields",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		globalCfg.InitConfig()
		qbutil.RequireTableID(globalCfg)

		input := &qb.GetSchemaInput{ID: globalCfg.TableID()}

		client := qb.NewClient(globalCfg)
		output, err := client.GetSchema(input)
		cliutil.HandleError(err, "error executing request")

		s, err := formatFields(output)
		cliutil.HandleError(err, "error formatting output")

		fmt.Println(s)
	},
}

func init() {
	fieldCmd.AddCommand(fieldListCmd)
	fieldListCfg = cliutil.InitConfig(qb.EnvVarPrefix)
}

func formatFields(out qb.GetSchemaOutput) (string, error) {

	// Build map of field ID to labels.
	fields := make(map[int]string)
	for _, f := range out.Fields {
		fields[f.FieldID] = f.Label
	}

	// Format and render the output.
	v := &FieldsRenderedJSON{Fields: fields}
	b, err := json.MarshalIndent(v, "", "    ")
	if err != nil {
		return "", err
	}

	return string(b), nil
}

// FieldsRenderedJSON renders records in JSON.
type FieldsRenderedJSON struct {
	Fields map[int]string `json:"fields"`
}
