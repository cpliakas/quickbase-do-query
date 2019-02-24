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
	"github.com/cpliakas/quickbase-do-query/qb"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var recordAddCfg *viper.Viper

var recordAddCmd = &cobra.Command{
	Use:   "add [FIELD_VALUES]",
	Short: "Adds a record",
	Long:  ``,
	Args:  recordAddCmdValidate,
	Run: func(cmd *cobra.Command, args []string) {

		values := cliutil.ParseKeyValue(strings.Join(args, " "))
		fields, err := parseValues(values)
		cliutil.HandleError(err, "error parsing field values")

		input := &qb.AddRecordInput{
			TableID: globalCfg.TableID(),
			Fields:  fields,
		}

		client := qb.NewClient(globalCfg)
		output, err := client.AddRecord(input)
		cliutil.HandleError(err, "error formatting output")

		cliutil.PrintJSON(output)
	},
}

func init() {
	recordCmd.AddCommand(recordAddCmd)
	recordAddCfg = cliutil.InitConfig(qb.EnvVarPrefix)
}

func recordAddCmdValidate(cmd *cobra.Command, args []string) error {
	globalCfg.RequireTableID = true
	if err := globalCfg.Validate(); err != nil {
		return err
	}

	if len(args) < 1 {
		return errors.New("missing required argument: [FIELD_VALUES]")
	}

	return nil
}

// parseValues parses the values argument into a qb.AddRecordInputField slice.
// TODO Make this generic?
func parseValues(m map[string]string) ([]qb.AddRecordInputField, error) {
	isNumeric := regexp.MustCompile(`^[1-9][0-9]*$`)
	fields := make([]qb.AddRecordInputField, len(m))

	i := 0
	for field, value := range m {

		// Add the field ID or field label.
		fields[i] = qb.AddRecordInputField{}
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
