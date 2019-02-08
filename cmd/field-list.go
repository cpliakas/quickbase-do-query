package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/cpliakas/quickbase-do-query/quickbase"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// TODO: Replace panics with a better error handling mechanism.
var fieldListCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists fields",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		cfg := quickbase.NewConfig()
		client := quickbase.NewClient(cfg)

		input := &quickbase.GetSchemaInput{ID: viper.GetString("table-id")}

		output, err := client.GetSchema(input)
		if err != nil {
			panic(err)
		}

		s, err := formatFields(output)
		if err != nil {
			panic(err)
		}

		fmt.Println(s)
	},
}

func init() {
	fieldCmd.AddCommand(fieldListCmd)
}

func formatFields(out quickbase.GetSchemaOutput) (string, error) {

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
