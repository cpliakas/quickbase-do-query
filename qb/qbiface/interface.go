package qbiface

import (
	"github.com/cpliakas/quickbase-do-query/qb"
)

// ClientAPI provides an interface to enable mocking the Quick Base service
// client's API calls.
type ClientAPI interface {
	AddRecord(*qb.AddRecordInput) (qb.AddRecordOutput, error)
	Authenticate(*qb.AuthenticateInput) (qb.AuthenticateOutput, error)
	DoQuery(*qb.DoQueryInput) (qb.DoQueryOutput, error)
	EditRecord(*qb.EditRecordInput) (qb.EditRecordOutput, error)
	GetSchema(*qb.GetSchemaInput) (qb.GetSchemaOutput, error)
	ImportFromCSV(*qb.ImportFromCSVInput) (qb.ImportFromCSVOutput, error)
	SetVariable(*qb.SetVariableInput) (qb.SetVariableOutput, error)
	UploadFile(*qb.UploadFileInput) (qb.UploadFileOutput, error)
}
