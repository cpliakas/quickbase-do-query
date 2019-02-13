package cmd

import (
	"encoding/base64"
	"errors"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/cpliakas/quickbase-do-query/cliutil"
	qb "github.com/cpliakas/quickbase-do-query/quickbase"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var recordEditCfg *viper.Viper

var recordEditCmd = &cobra.Command{
	Use:   "edit [FIELD_VALUES]",
	Short: "Edits a record",
	Long:  ``,
	Args:  recordEditCmdValidate,
	Run: func(cmd *cobra.Command, args []string) {

		values := cliutil.ParseKeyValue(strings.Join(args, " "))
		fields, err := parseEditValues(values)
		cliutil.HandleError(err, "error parsing field values")

		input := &qb.EditRecordInput{
			TableID:  globalCfg.TableID(),
			RecordID: recordEditCfg.GetInt("record-id"),
			Fields:   fields,
		}

		client := qb.NewClient(globalCfg)
		output, err := client.EditRecord(input)
		cliutil.HandleError(err, "error formatting output")

		cliutil.PrintJSON(output)
	},
}

func init() {
	recordCmd.AddCommand(recordEditCmd)
	recordEditCfg = cliutil.InitConfig(qb.EnvVarPrefix)

	flags := cliutil.NewFlagger(recordEditCmd, recordEditCfg)
	flags.Int("record-id", "r", 0, "ID of the record being edited")
}

func recordEditCmdValidate(cmd *cobra.Command, args []string) error {
	globalCfg.RequireTableID = true
	if err := globalCfg.Validate(); err != nil {
		return err
	}

	if len(args) < 1 {
		return errors.New("missing required argument: [FIELD_VALUES]")
	}
	if recordEditCfg.GetInt("record-id") <= 0 {
		return errors.New("missing required option: record-id")
	}

	return nil
}

// parseValues parses the values argument into a qb.EditRecordInputField slice.
// TODO Make this generic?
func parseEditValues(m map[string]string) ([]qb.EditRecordInputField, error) {
	isNumeric := regexp.MustCompile(`^[1-9][0-9]*$`)
	fields := make([]qb.EditRecordInputField, len(m))

	i := 0
	for field, value := range m {

		// Add the field ID or field label.
		fields[i] = qb.EditRecordInputField{}
		if isNumeric.MatchString(field) {
			fid, _ := strconv.Atoi(field)
			fields[i].ID = fid
		} else {
			fields[i].Label = field
		}

		// Add the value, extract data from a file if appropriate.
		if !strings.HasPrefix(value, "file://") {
			fields[i].Value = value
		} else {
			filePath := strings.TrimPrefix(value, "file://")
			fileData, err := ioutil.ReadFile(filePath)
			if err != nil {
				return fields, err
			}
			fields[i].Value = base64.StdEncoding.EncodeToString(fileData)
			fields[i].FileName = filepath.Base(filePath)
		}

		i = i + 1
	}

	return fields, nil
}
