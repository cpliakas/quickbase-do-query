package cmd

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"

	"github.com/cpliakas/quickbase-do-query/cliutil"
	"github.com/cpliakas/quickbase-do-query/qbutil"
	qb "github.com/cpliakas/quickbase-do-query/quickbase"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var fileAttachmentUploadCfg *viper.Viper

var fileAttachmentUploadCmd = &cobra.Command{
	Use:     "upload [RECORD_ID] [FIELD_ID] [FILEPATH]",
	Short:   "Uploads a file attachment",
	Long:    ``,
	PreRunE: globalCfg.PreRunE,
	Run: func(cmd *cobra.Command, args []string) {
		qbutil.RequireTableID(globalCfg)

		rid := cliutil.RequireArgInt(args, 0, "record-id")
		fid := cliutil.RequireArgInt(args, 1, "field-id")
		filepath := cliutil.RequireArg(args, 2, "filepath")

		fileData, err := ioutil.ReadFile(filepath)
		cliutil.HandleError(err, "error reading file")

		field := qb.UploadFileAttachmentInputField{
			FieldID:  fid,
			FileData: base64.StdEncoding.EncodeToString(fileData),
			Name:     "test.txt",
		}

		input := &qb.UploadFileAttachmentInput{
			TableID:  globalCfg.TableID(),
			RecordID: rid,
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
}
