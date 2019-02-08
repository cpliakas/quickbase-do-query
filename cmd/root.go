package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/cpliakas/quickbase-do-query/quickbase"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// TODO: Replace panics with ba better error handling mechanism.
var rootCmd = &cobra.Command{
	Use:   "quickbase-do-query",
	Short: "A command line tool that gets records from a Quick Base table.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		cfg := quickbase.NewConfig()
		client := quickbase.NewClient(cfg)

		input := &quickbase.DoQueryInput{}
		input.TableID = viper.GetString("table-id")

		query := viper.GetString("query")
		queryID := viper.GetInt("query-id")
		queryName := viper.GetString("query-name")

		if query != "" {
			input.Query = query
		} else if queryID > 0 {
			input.QueryID = queryID
		} else if queryName != "" {
			input.QueryName = queryName
		}

		// TODO: Support these options.
		// Unsorted()
		// OnlyNew()
		// ReturnPercentage

		fields, err := parseFieldsOption(viper.GetString("fields"))
		if err != nil {
			panic(err)
		}
		input.FieldSlice = fields

		input.Offset(viper.GetInt("offset"))
		input.Limit(viper.GetInt("limit"))

		sort, order, err := parseSortOption(viper.GetString("sort"))
		if err != nil {
			panic(err)
		}
		input.SortSlice = sort
		input.EnsureOptions().SortOrderSlice = order

		output, err := client.DoQuery(input)
		if err != nil {
			panic(err)
		}

		s, err := formatRecords(output, false) // @TODO Hijack the UseFIDs in Input.
		if err != nil {
			panic(err)
		}

		fmt.Println(s)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {

	viper.SetEnvPrefix("QUICKBASE")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	rootCmd.PersistentFlags().StringP("realm-host", "R", "", "The realm host")
	viper.BindPFlag("realm-host", rootCmd.PersistentFlags().Lookup("realm-host"))

	rootCmd.PersistentFlags().StringP("app-id", "I", "", "The application ID")
	viper.BindPFlag("app-id", rootCmd.PersistentFlags().Lookup("app-id"))

	rootCmd.PersistentFlags().StringP("user-token", "U", "", "The user token used to authenticate the request")
	viper.BindPFlag("user-token", rootCmd.PersistentFlags().Lookup("user-token"))

	rootCmd.PersistentFlags().StringP("ticket", "T", "", "The ticket used to authenticate the request")
	viper.BindPFlag("ticket", rootCmd.PersistentFlags().Lookup("ticket"))

	rootCmd.PersistentFlags().StringP("app-token", "A", "", "The app token used to authenticate the request")
	viper.BindPFlag("app-token", rootCmd.PersistentFlags().Lookup("app-token"))

	rootCmd.PersistentFlags().StringP("table-id", "t", "", "The table to get records from")
	viper.BindPFlag("table-id", rootCmd.PersistentFlags().Lookup("table-id"))

	rootCmd.Flags().StringP("query", "q", "", "The query to get records from the table")
	viper.BindPFlag("query", rootCmd.Flags().Lookup("query"))

	rootCmd.Flags().StringP("query-id", "i", "", "The ID of the query that gets records from the table")
	viper.BindPFlag("query-id", rootCmd.Flags().Lookup("query-id"))

	rootCmd.Flags().StringP("query-name", "n", "", "The name of the query that gets records from the table")
	viper.BindPFlag("query-name", rootCmd.Flags().Lookup("query-name"))

	rootCmd.Flags().StringP("fields", "f", "", "The fields to return from the table")
	viper.BindPFlag("fields", rootCmd.Flags().Lookup("fields"))

	rootCmd.Flags().IntP("limit", "l", 25, "The maximum number of records to return")
	viper.BindPFlag("limit", rootCmd.Flags().Lookup("limit"))

	rootCmd.Flags().IntP("offset", "o", 0, "The number of records to skip")
	viper.BindPFlag("offset", rootCmd.Flags().Lookup("offset"))

	rootCmd.Flags().StringP("sort", "s", "", "A comma-delimited list of fields to sort by")
	viper.BindPFlag("sort", rootCmd.Flags().Lookup("sort"))

	viper.ReadInConfig()
}

// formatRecords formats records in JSON.
func formatRecords(out quickbase.DoQueryOutput, useFIDs bool) (string, error) {

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

// parseFieldsOption parses the sort option passed through the command line.
func parseFieldsOption(fieldsStr string) ([]int, error) {
	if fieldsStr == "" {
		return []int{}, nil
	}

	// TODO: Replace wth regex to allow for spaces after comma.
	parts := strings.Split(fieldsStr, ",")
	fields := make([]int, len(parts))

	for k, part := range parts {
		fid, err := strconv.Atoi(part)
		if err != nil {
			// TODO: Invalid input error instead of generic.
			return []int{}, errors.New("invalid field ID")
		}
		fields[k] = fid
	}

	return fields, nil
}

// parseSortOption parses the sort option passed through the command line.
//
// The first two parameters are fids to sort on and flow (e.g. A,D) respectively.
func parseSortOption(sortStr string) ([]int, []string, error) {
	if sortStr == "" {
		return []int{}, []string{}, nil
	}

	// TODO: Replace wth regex to allow for spaces after comma.
	parts := strings.Split(sortStr, ",")
	sort := make([]int, len(parts))
	order := make([]string, len(parts))

	re := regexp.MustCompile(`^([0-9]+)\s*(D|A|DESC|ASC)?$`)
	for k, part := range parts {

		match := re.FindStringSubmatch(part)
		if len(match) == 0 {
			// TODO: Invalid input error instead of generic.
			return []int{}, []string{}, errors.New("invalid input")
		}

		fid, err := strconv.Atoi(match[1])
		if err != nil {
			// TODO: Invalid input error instead of generic.
			return []int{}, []string{}, errors.New("invalid field ID")
		}
		sort[k] = fid

		// TODO: Validate whether match[2] exists?
		order[k] = match[2]
		if order[k] == "DESC" || order[k] == "D" {
			order[k] = "D"
		} else {
			order[k] = "A"
		}
	}

	return sort, order, nil
}
