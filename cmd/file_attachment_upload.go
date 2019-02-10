package cmd

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/cpliakas/quickbase-do-query/cliutil"
	qb "github.com/cpliakas/quickbase-do-query/quickbase"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var fileAttachmentUploadCfg *viper.Viper

var fileAttachmentUploadCmd = &cobra.Command{
	Use:   "upload [FILEPATH]",
	Short: "Uploads a file attachment",
	Long:  ``,
	Args:  fileAttachmentUploadCmdValidate,
	Run: func(cmd *cobra.Command, args []string) {

		fileData, err := ioutil.ReadFile(args[0])
		cliutil.HandleError(err, "error reading file")

		fileName := fileAttachmentUploadCfg.GetString("file-name")
		if fileName == "" {
			fileName = filepath.Base(args[0])
		}

		field := qb.UploadFileAttachmentInputField{
			FieldID:  fileAttachmentUploadCfg.GetInt("field-id"),
			FileData: base64.StdEncoding.EncodeToString(fileData),
			Name:     fileName,
		}

		input := &qb.UploadFileAttachmentInput{
			TableID:  globalCfg.TableID(),
			RecordID: fileAttachmentUploadCfg.GetInt("record-id"),
			Fields:   []qb.UploadFileAttachmentInputField{field},
		}

		client := qb.NewClient(globalCfg)
		output, err := client.UploadFileAttachment(input)
		cliutil.HandleError(err, "error formatting output")

		fmt.Println(output.Fields[0].URL)
	},
}

func init() {
	fileAttachmentCmd.AddCommand(fileAttachmentUploadCmd)
	fileAttachmentUploadCfg = cliutil.InitConfig(qb.EnvVarPrefix)

	flags := cliutil.NewFlagger(fileAttachmentUploadCmd, fileAttachmentUploadCfg)
	flags.Int("field-id", "f", 0, "the file attachment's field ID")
	flags.String("file-name", "n", "", "the name of file stored in the record")
	flags.Int("record-id", "r", 0, "record ID the file attachment is being uploaded to")
}

func fileAttachmentUploadCmdValidate(cmd *cobra.Command, args []string) error {
	globalCfg.RequireTableID = true
	if err := globalCfg.Validate(); err != nil {
		return err
	}

	if len(args) < 1 {
		return errors.New("missing required argument: [FILEPATH]")
	}

	if fileAttachmentUploadCfg.GetInt("field-id") <= 0 {
		return errors.New("missing required option: field-id")
	}
	if fileAttachmentUploadCfg.GetInt("record-id") <= 0 {
		return errors.New("missing required option: record-id")
	}

	return nil
}
