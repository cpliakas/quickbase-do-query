package cmd

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/cpliakas/quickbase-do-query/cliutil"
	"github.com/cpliakas/quickbase-do-query/qbutil"
	qb "github.com/cpliakas/quickbase-do-query/quickbase"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var doQueryCfg *viper.Viper

var doQueryCmd = &cobra.Command{
	Use:   "do-query",
	Short: "Executes a query.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		globalCfg.InitConfig()
		qbutil.RequireTableID(globalCfg)

		input := &qb.DoQueryInput{}
		input.TableID = globalCfg.TableID()

		query := doQueryCfg.GetString("query")
		queryID := doQueryCfg.GetInt("query-id")
		queryName := doQueryCfg.GetString("query-name")

		if query != "" {
			input.Query = query
		} else if queryID > 0 {
			input.QueryID = queryID
		} else if queryName != "" {
			input.QueryName = queryName
		} else {
			err := errors.New("query, query-id, or query-name")
			cliutil.HandleError(err, "missing required option")
		}

		// TODO: Support these options.
		// Unsorted()
		// OnlyNew()
		// ReturnPercentage

		fields, err := qbutil.ParseFieldsOption(doQueryCfg.GetString("fields"))
		cliutil.HandleError(err, "fields option invalid")
		input.FieldSlice = fields

		sort, order, err := qbutil.ParseSortOption(doQueryCfg.GetString("sort"))
		cliutil.HandleError(err, "sort option invalid")
		input.Sort(sort, order)

		input.Offset(doQueryCfg.GetInt("offset"))
		input.Limit(doQueryCfg.GetInt("limit"))

		client := qb.NewClient(globalCfg)
		output, err := client.DoQuery(input)
		cliutil.HandleError(err, "error executing request")

		// @TODO Hijack the UseFIDs in Input.
		s, err := formatRecords(output, false)
		cliutil.HandleError(err, "error formatting output")

		fmt.Println(s)
	},
}

func init() {
	rootCmd.AddCommand(doQueryCmd)

	doQueryCfg = cliutil.InitConfig(qb.EnvVarPrefix)
	flags := cliutil.NewFlagger(doQueryCmd, doQueryCfg)

	flags.String("query", "q", "", "query that gets records from the table")
	flags.String("query-id", "i", "", "ID of the query that gets records from the table")
	flags.String("query-name", "n", "", "name of the query that gets records from the table")
	flags.String("fields", "f", "", "comma-delimited list of fields to return")
	flags.Int("limit", "l", 25, "maximum number of records to return")
	flags.Int("offset", "o", 0, "number of records to skip")
	flags.String("sort", "s", "", "comma-delimited list of fields to sort by")
}

// formatRecords formats records in JSON.
func formatRecords(out qb.DoQueryOutput, useFIDs bool) (string, error) {

	// Build a field map so we can key the field by label.
	fieldMap := make(map[int]string)
	for _, f := range out.Fields {
		fieldMap[f.FieldID] = f.Label
	}

	// Builds the rendered output.
	records := make([]RecordRenderedJSON, len(out.Records))
	for k, r := range out.Records {
		records[k].RecordID = r.RecordID
		records[k].UpdateID = r.UpdateID
		records[k].Fields = make(map[string]interface{})
		for _, f := range r.Fields {
			label := fieldMap[f.FieldID]
			records[k].Fields[label] = f.Value
		}
	}

	// Format and render the output.
	v := &RecordsRenderedJSON{Records: records}
	b, err := json.MarshalIndent(v, "", "    ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// RecordsRenderedJSON renders records in JSON.
type RecordsRenderedJSON struct {
	Records []RecordRenderedJSON `json:"records"`
}

// RecordRenderedJSON a record in JSON.
type RecordRenderedJSON struct {
	RecordID int                    `json:"record-id"`
	UpdateID int                    `json:"update-id"`
	Fields   map[string]interface{} `json:"fields"`
}
