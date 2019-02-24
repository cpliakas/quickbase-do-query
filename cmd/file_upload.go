package cmd

import (
	"encoding/base64"
	"errors"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/cpliakas/quickbase-do-query/cliutil"
	"github.com/cpliakas/quickbase-do-query/qb"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var fileUploadCfg *viper.Viper

var fileUploadCmd = &cobra.Command{
	Use:   "upload [FILEPATH]",
	Short: "Uploads a file ",
	Long:  ``,
	Args:  fileUploadCmdValidate,
	Run: func(cmd *cobra.Command, args []string) {

		filePath := strings.TrimPrefix(args[0], "file://")
		fileData, err := ioutil.ReadFile(filePath)
		cliutil.HandleError(err, "error reading file")

		fileName := fileUploadCfg.GetString("file-name")
		if fileName == "" {
			fileName = filepath.Base(filePath)
		}

		field := qb.UploadFileInputField{
			ID:       fileUploadCfg.GetInt("field-id"),
			FileData: base64.StdEncoding.EncodeToString(fileData),
			Name:     fileName,
		}

		input := &qb.UploadFileInput{
			TableID:  globalCfg.TableID(),
			RecordID: fileUploadCfg.GetInt("record-id"),
			Fields:   []qb.UploadFileInputField{field},
		}

		client := qb.NewClient(globalCfg)
		output, err := client.UploadFile(input)
		cliutil.HandleError(err, "error formatting output")

		cliutil.PrintJSON(output)
	},
}

func init() {
	fileCmd.AddCommand(fileUploadCmd)
	fileUploadCfg = cliutil.InitConfig(qb.EnvVarPrefix)

	flags := cliutil.NewFlagger(fileUploadCmd, fileUploadCfg)
	flags.Int("field-id", "f", 0, "the file 's field ID")
	flags.String("file-name", "n", "", "the name of file stored in the record")
	flags.Int("record-id", "r", 0, "record ID the file is being uploaded to")
}

func fileUploadCmdValidate(cmd *cobra.Command, args []string) error {
	globalCfg.RequireTableID = true
	if err := globalCfg.Validate(); err != nil {
		return err
	}

	if len(args) < 1 {
		return errors.New("missing required argument: [FILEPATH]")
	}

	if fileUploadCfg.GetInt("field-id") <= 0 {
		return errors.New("missing required option: field-id")
	}
	if fileUploadCfg.GetInt("record-id") <= 0 {
		return errors.New("missing required option: record-id")
	}

	return nil
}
