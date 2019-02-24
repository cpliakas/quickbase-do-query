package cmd

import (
	"errors"
	"strconv"

	"github.com/cpliakas/quickbase-do-query/cliutil"
	"github.com/cpliakas/quickbase-do-query/qb"
	"github.com/cpliakas/quickbase-do-query/qbutil"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var doQueryCfg *viper.Viper

var doQueryCmd = &cobra.Command{
	Use:   "query",
	Short: "Executes a query against a table",
	Long:  ``,
	Args:  doQueryCmdValidate,
	Run: func(cmd *cobra.Command, args []string) {
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
		// ReturnPercentage()

		fields, err := qbutil.ParseFieldsOption(doQueryCfg.GetString("fields"))
		cliutil.HandleError(err, "fields option invalid")
		input.FieldList = fields

		sort, order, err := qbutil.ParseSortOption(doQueryCfg.GetString("sort"))
		cliutil.HandleError(err, "sort option invalid")
		input.Sort(sort, order)

		input.Offset(doQueryCfg.GetInt("offset"))
		input.Limit(doQueryCfg.GetInt("limit"))

		client := qb.NewClient(globalCfg)
		output, err := client.DoQuery(input)
		cliutil.HandleError(err, "error executing request")

		v := newDoQueryOutput(output, doQueryCfg.GetBool("use-labels"))
		cliutil.PrintJSON(v)
	},
}

func init() {
	rootCmd.AddCommand(doQueryCmd)
	doQueryCfg = cliutil.InitConfig(qb.EnvVarPrefix)

	flags := cliutil.NewFlagger(doQueryCmd, doQueryCfg)
	flags.String("fields", "f", "", "comma-delimited list of fields to return")
	flags.Int("limit", "l", 25, "maximum number of records to return")
	flags.Int("offset", "o", 0, "number of records to skip")
	flags.String("query", "q", "", "query that gets records from the table")
	flags.String("query-id", "i", "", "ID of the query that gets records from the table")
	flags.String("query-name", "n", "", "name of the query that gets records from the table")
	flags.String("sort", "s", "", "comma-delimited list of fields to sort by")
	flags.Bool("use-labels", "u", false, "key by label instead of field ID")
}

func doQueryCmdValidate(cmd *cobra.Command, args []string) error {
	globalCfg.RequireTableID = true
	return globalCfg.Validate()
}

// newDoQueryOutput returns a DoQueryOutput.
func newDoQueryOutput(out qb.DoQueryOutput, useLabels bool) DoQueryOutput {

	// Build a field map so we can key the field by label.
	// TODO: Don't build a map
	fieldMap := make(map[int]string)
	for _, f := range out.Fields {
		fieldMap[f.FieldID] = f.Label
	}

	// Builds the rendered output.
	records := make([]DoQueryOutputRecord, len(out.Records))
	for k, r := range out.Records {

		records[k].ID = r.RecordID
		records[k].UpdateID = r.UpdateID
		records[k].Fields = make(map[string]interface{})

		for _, f := range r.Fields {
			var label string
			if !useLabels {
				label = strconv.Itoa(f.FieldID)
			} else {
				label = fieldMap[f.FieldID]
			}
			records[k].Fields[label] = f.Value
		}
	}

	return DoQueryOutput{
		UserData: out.UserData,
		Records:  records,
	}
}

// DoQueryOutput models the output that prints the matched records.
type DoQueryOutput struct {
	UserData string                `json:"user_data,omitempty"`
	Records  []DoQueryOutputRecord `json:"records"`
}

// DoQueryOutputRecord models the output that prints a record.
type DoQueryOutputRecord struct {
	ID       int                    `json:"record_id"`
	UpdateID int                    `json:"update_id"`
	Fields   map[string]interface{} `json:"fields"`
}
