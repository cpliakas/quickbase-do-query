package cmd

import (
	"errors"
	"io/ioutil"

	"github.com/cpliakas/quickbase-do-query/cliutil"
	"github.com/cpliakas/quickbase-do-query/qbutil"
	qb "github.com/cpliakas/quickbase-do-query/quickbase"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var csvImportCfg *viper.Viper

var csvImportCmd = &cobra.Command{
	Use:   "import [FILEPATH]",
	Short: "imports data from a CSV file into a table",
	Long:  ``,
	Args:  csvImportCmdValidate,
	Run: func(cmd *cobra.Command, args []string) {

		fields, err := qbutil.ParseFieldsOption(csvImportCfg.GetString("fields"))
		cliutil.HandleError(err, "fields option invalid")

		input := &qb.ImportFromCSVInput{
			TableID:          globalCfg.TableID(),
			FieldList:        fields,
			SkipFirstRow:     qb.Bool(csvImportCfg.GetBool("skip-first-row")),
			DecimalAsPercent: qb.Bool(csvImportCfg.GetBool("decimal-as-percent")),
		}

		if mfid := csvImportCfg.GetInt("merge-field-id"); mfid > 0 {
			input.MergeFieldID = mfid
		}

		fileData, err := ioutil.ReadFile(args[0])
		cliutil.HandleError(err, "error reading file")
		input.CSV(fileData)

		client := qb.NewClient(globalCfg)
		output, err := client.ImportFromCSV(input)
		cliutil.HandleError(err, "error formatting output")

		// TODO Nice output
		cliutil.PrintJSON(output)
	},
}

func init() {
	csvCmd.AddCommand(csvImportCmd)
	csvImportCfg = cliutil.InitConfig(qb.EnvVarPrefix)

	flags := cliutil.NewFlagger(csvImportCmd, csvImportCfg)
	flags.Bool("decimal-as-percent", "d", false, "decimal values like 0.50 sre interpreted to mean 50%")
	flags.String("fields", "f", "", "list of fields to import")
	flags.Int("merge-field-id", "m", 0, "use as the key field")
	flags.String("output-fields", "o", "", "list of fields to return")
	flags.Bool("skip-first-row", "s", false, "do not importing the first row of data")
}

func csvImportCmdValidate(cmd *cobra.Command, args []string) error {
	globalCfg.RequireTableID = true
	if err := globalCfg.Validate(); err != nil {
		return err
	}

	if len(args) < 1 {
		return errors.New("missing required argument: [FILEPATH]")
	}

	return nil
}
