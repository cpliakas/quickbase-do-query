package quickbase

import (
	"bytes"
	"encoding/csv"
	"encoding/xml"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// AuthenticateInput models requests sent to API_Authenticate.
// See https://help.quickbase.com/api-guide/authenticate.html
type AuthenticateInput struct {
	RequestParams

	Hours    int    `xml:"hours,omitempty"`
	Password string `xml:"password"`
	Username string `xml:"username"`
}

func (input *AuthenticateInput) method() string           { return http.MethodPost }
func (input *AuthenticateInput) uri() string              { return "/db/main" }
func (input *AuthenticateInput) payload() ([]byte, error) { return xml.Marshal(input) }
func (input *AuthenticateInput) headers(req *http.Request) {
	req.Header.Set("Content-Type", "application/xml")
	req.Header.Set("QUICKBASE-ACTION", "API_Authenticate")
}

// AuthenticateOutput models the response returned by API_Authenticate.
// See https://help.quickbase.com/api-guide/authenticate.html
type AuthenticateOutput struct {
	ResponseParams

	Ticket string `xml:"ticket"`
	UserID string `xml:"userid"`
}

func (output *AuthenticateOutput) parse(body []byte, res *http.Response) error {
	return parseXML(output, body, res)
}

// Authenticate makes call to API_Authenticate.
// See https://help.quickbase.com/api-guide/authenticate.html
func (c Client) Authenticate(input *AuthenticateInput) (output AuthenticateOutput, err error) {
	err = c.Do(input, &output)
	if err == nil && output.ErrorCode != 0 {
		err = fmt.Errorf("error executing API_Authenticate: %s (error code: %v)", output.ErrorText, output.ErrorCode)
	}
	return
}

// DoQueryInput models the request sent to API_DoQuery.
// See https://help.quickbase.com/api-guide/do_query.html
type DoQueryInput struct {
	RequestParams
	Credentials

	TableID          string               `xml:"-"`
	Query            string               `xml:"query,omitempty"`
	QueryID          int                  `xml:"qid,omitempty"`
	QueryName        string               `xml:"qname,omitempty"`
	IncludeRecordIDs Bool                 `xml:"includeRids,omitempty"`
	ReturnPercentage Bool                 `xml:"returnpercentage,omitempty"`
	UseFieldIDs      Bool                 `xml:"useFids,omitempty"`
	Format           string               `xml:"fmt,omitempty"`
	FieldList        FieldList            `xml:"clist,omitempty"`
	SortList         FieldList            `xml:"slist,omitempty"`
	Options          *DoQueryInputOptions `xml:"options,omitempty"`
}

func (input *DoQueryInput) setCredentials(creds Credentials) { input.Credentials = creds }
func (input *DoQueryInput) method() string                   { return http.MethodPost }
func (input *DoQueryInput) uri() string                      { return "/db/" + input.TableID }
func (input *DoQueryInput) payload() ([]byte, error)         { return xml.Marshal(input) }
func (input *DoQueryInput) headers(req *http.Request) {
	req.Header.Set("Content-Type", "application/xml")
	req.Header.Set("QUICKBASE-ACTION", "API_DoQuery")
}

// EnsureOptions returns an initialized Options property. This method should be
// used in favor of accessing the property directly to avoid null pointer
// exceptions.
func (input *DoQueryInput) EnsureOptions() *DoQueryInputOptions {
	if input.Options == nil {
		input.Options = &DoQueryInputOptions{}
	}
	return input.Options
}

// Fields sets the fields that are returned in the response.
func (input *DoQueryInput) Fields(fids ...int) *DoQueryInput {
	input.FieldList = fids
	return input
}

// SortBy sets the fields to be sorted by.
func (input *DoQueryInput) SortBy(fids ...int) *DoQueryInput {
	input.SortList = fids
	return input
}

// SortOrder sets the "sortorder" option.
func (input *DoQueryInput) SortOrder(order ...string) *DoQueryInput {
	input.EnsureOptions().SortOrderList = order
	return input
}

// Sort sets the fields and order in one shot.
func (input *DoQueryInput) Sort(sort []int, order []string) *DoQueryInput {
	input.SortList = sort
	input.EnsureOptions().SortOrderList = order
	return input
}

// Limit sets the "num" option.
func (input *DoQueryInput) Limit(n int) *DoQueryInput {
	input.EnsureOptions().Limit = n
	return input
}

// Offset sets the "skp" option.
func (input *DoQueryInput) Offset(n int) *DoQueryInput {
	input.EnsureOptions().Offset = n
	return input
}

// OnlyNew sets the "onlynew" option.
func (input *DoQueryInput) OnlyNew() *DoQueryInput {
	input.EnsureOptions().OnlyNew = true
	return input
}

// Unsorted sets the "nosort" option.
func (input *DoQueryInput) Unsorted() *DoQueryInput {
	input.EnsureOptions().Unsorted = true
	return input
}

// DoQueryInputOptions models the "options" element in API_DoQuery requests.
type DoQueryInputOptions struct {
	SortOrderList []string
	Limit         int
	Offset        int
	OnlyNew       bool
	Unsorted      bool
}

// MarshalXML implements Marshaler.MarshalXML and formats the value of the
// "options" element.
func (o DoQueryInputOptions) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	opts := []string{}

	if o.Offset > 0 {
		opts = append(opts, "skp-"+strconv.Itoa(o.Offset))
	}
	if o.Limit > 0 {
		opts = append(opts, "num-"+strconv.Itoa(o.Limit))
	}
	if len(o.SortOrderList) > 0 {
		opts = append(opts, "sortorder-"+strings.Join(o.SortOrderList, ""))
	}
	if o.OnlyNew == true {
		opts = append(opts, "onlynew")
	}
	if o.Unsorted == true {
		opts = append(opts, "nosort")
	}

	return e.EncodeElement(strings.Join(opts, "."), start)
}

// DoQueryOutput models the response returned by API_DoQuery.
// See https://help.quickbase.com/api-guide/do_query.html.
type DoQueryOutput struct {
	ResponseParams

	Fields  []DoQueryOutputField  `xml:"table>fields>field"`
	Records []DoQueryOutputRecord `xml:"table>records>record"`
}

// DoQueryOutputField models the "table>fields" element. in an API_DoQuery
// response.
type DoQueryOutputField struct {
	FieldID  int    `xml:"id,attr"`
	Type     string `xml:"field_type,attr"`
	BaseType string `xml:"base_type,attr"`
	Mode     string `xml:"mode,attr"`
	Label    string `xml:"label"`
}

// DoQueryOutputRecord models the "table>records>record" element in an
// API_DoQuery response.
type DoQueryOutputRecord struct {
	RecordID int                        `xml:"rid,attr"`
	UpdateID int                        `xml:"update_id"`
	Fields   []DoQueryOutputRecordField `xml:"f"`
}

// DoQueryOutputRecordField models the "table>records>record>field" element
// in an API_DoQuery response.
type DoQueryOutputRecordField struct {
	FieldID int    `xml:"id,attr"`
	Value   string `xml:",chardata"`
}

func (output *DoQueryOutput) parse(body []byte, res *http.Response) error {
	return parseXML(output, body, res)
}

// DoQuery makes call to API_DoQuery.
// See https://help.quickbase.com/api-guide/do_query.html.
func (c Client) DoQuery(input *DoQueryInput) (output DoQueryOutput, err error) {
	// Required for predictable output.
	input.Format = "structured"
	input.IncludeRecordIDs = true

	err = c.Do(input, &output)
	if err == nil && output.ErrorCode != 0 {
		err = fmt.Errorf("error executing API_DoQuery: %s (error code: %v)", output.ErrorText, output.ErrorCode)
	}
	return
}

// GetSchemaInput models requests sent to API_GetSchema.
// See https://help.quickbase.com/api-guide/getschema.html
type GetSchemaInput struct {
	RequestParams
	Credentials

	ID string `xml:"-"`
}

func (input *GetSchemaInput) setCredentials(creds Credentials) { input.Credentials = creds }
func (input *GetSchemaInput) method() string                   { return http.MethodPost }
func (input *GetSchemaInput) uri() string                      { return "/db/" + input.ID }
func (input *GetSchemaInput) payload() ([]byte, error)         { return xml.Marshal(input) }
func (input *GetSchemaInput) headers(req *http.Request) {
	req.Header.Set("Content-Type", "application/xml")
	req.Header.Set("QUICKBASE-ACTION", "API_GetSchema")
}

// GetSchemaOutput models responses returned API_GetSchema.
// See https://help.quickbase.com/api-guide/getschema.html
type GetSchemaOutput struct {
	ResponseParams

	Fields []DoQueryOutputField `xml:"table>fields>field"`
}

func (output *GetSchemaOutput) parse(body []byte, res *http.Response) error {
	return parseXML(output, body, res)
}

// GetSchema makes call to API_GetSchema.
// See https://help.quickbase.com/api-guide/getschema.html
func (c Client) GetSchema(input *GetSchemaInput) (output GetSchemaOutput, err error) {
	err = c.Do(input, &output)
	if err == nil && output.ErrorCode != 0 {
		err = fmt.Errorf("error executing API_GetSchema: %s (error code: %v)", output.ErrorText, output.ErrorCode)
	}
	return
}

// ImportFromCSVInput models the request sent to API_ImportFromCSV
// See https://help.quickbase.com/api-guide/importfromcsv.html
type ImportFromCSVInput struct {
	RequestParams
	Credentials

	TableID          string                     `xml:"-"`
	FieldList        FieldList                  `xml:"clist,omitempty"`
	OutputFieldList  FieldList                  `xml:"clist_output,omitempty"`
	MergeFieldID     int                        `xml:"mergeFieldId,omitempty"`
	DecimalAsPercent Bool                       `xml:"decimalPercent,omitempty"`
	Records          *ImportFromCSVInputRecords `xml:"records_CSV"`
	SkipFirstRow     Bool                       `xml:"skipfirst,omitempty"`
}

func (input *ImportFromCSVInput) setCredentials(creds Credentials) { input.Credentials = creds }
func (input *ImportFromCSVInput) method() string                   { return http.MethodPost }
func (input *ImportFromCSVInput) uri() string                      { return "/db/" + input.TableID }
func (input *ImportFromCSVInput) payload() ([]byte, error)         { return xml.Marshal(input) }
func (input *ImportFromCSVInput) headers(req *http.Request) {
	req.Header.Set("Content-Type", "application/xml")
	req.Header.Set("QUICKBASE-ACTION", "API_ImportFromCSV")
}

// EnsureRecords returns an initialized Records property. This method should be
// used in favor of accessing the property directly to avoid null pointer
// exceptions.
func (input *ImportFromCSVInput) EnsureRecords() *ImportFromCSVInputRecords {
	if input.Records == nil {
		input.Records = &ImportFromCSVInputRecords{}
	}
	return input.Records
}

// CSV sets raw CSV data,
func (input *ImportFromCSVInput) CSV(csv []byte) {
	input.EnsureRecords().CSV = string(csv)
}

// FormatCSV converts a string slice slice ([][]string) into CSV data.
func (input *ImportFromCSVInput) FormatCSV(records [][]string) {
	var buf bytes.Buffer
	w := csv.NewWriter(&buf)
	w.WriteAll(records)
	input.EnsureRecords().CSV = buf.String()
}

// ImportFromCSVInputRecords models the "records_CSV" element in
// API_ImportFromCSV requests.
type ImportFromCSVInputRecords struct {
	CSV string `xml:",cdata"`
}

// ImportFromCSVOutput models the response returned by API_ImportFromCSV
// See https://help.quickbase.com/api-guide/importfromcsv.html
type ImportFromCSVOutput struct {
	ResponseParams

	NumRecordsAdded   int                         `xml:"num_recs_added"`
	NumRecordsInput   int                         `xml:"num_recs_input"`
	NumRecordsUpdated int                         `xml:"num_recs_updated"`
	Records           []ImportFromCSVOutputRecord `xml:"rids>rid,omitempty"`
}

func (output *ImportFromCSVOutput) parse(body []byte, res *http.Response) error {
	return parseXML(output, body, res)
}

// ImportFromCSVOutputRecord models the "rids>rid" element in API_ImportFromCSV
// responses.
type ImportFromCSVOutputRecord struct {
	ID       int `xml:",chardata"`
	UpdateID int `xml:"update_id,attr"`
}

// ImportFromCSV makes an API_ImportFromCSV call.
// See https://help.quickbase.com/api-guide/importfromcsv.html
func (c Client) ImportFromCSV(input *ImportFromCSVInput) (output ImportFromCSVOutput, err error) {
	err = c.Do(input, &output)
	if err == nil && output.ErrorCode != 0 {
		err = fmt.Errorf("error executing API_ImportFromCSV: %s (error code: %v)", output.ErrorText, output.ErrorCode)
	}
	return
}

// SetVariableInput models the request sent to API_SetDBvar
// See https://help.quickbase.com/api-guide/setdbvar.html
type SetVariableInput struct {
	RequestParams
	Credentials

	AppID string `xml:"-"`
	Name  string `xml:"varname"`
	Value string `xml:"value"`
}

func (input *SetVariableInput) setCredentials(creds Credentials) { input.Credentials = creds }
func (input *SetVariableInput) method() string                   { return http.MethodPost }
func (input *SetVariableInput) uri() string                      { return "/db/" + input.AppID }
func (input *SetVariableInput) payload() ([]byte, error)         { return xml.Marshal(input) }
func (input *SetVariableInput) headers(req *http.Request) {
	req.Header.Set("Content-Type", "application/xml")
	req.Header.Set("QUICKBASE-ACTION", "API_SetDBvar")
}

// SetVariableOutput models the response returned by API_SetDBvar
// See https://help.quickbase.com/api-guide/setdbvar.html
type SetVariableOutput struct {
	ResponseParams
}

func (output *SetVariableOutput) parse(body []byte, res *http.Response) error {
	return parseXML(output, body, res)
}

// SetVariable makes an API_SetDBvar call.
// See https://help.quickbase.com/api-guide/setdbvar.html
func (c Client) SetVariable(input *SetVariableInput) (output SetVariableOutput, err error) {
	err = c.Do(input, &output)
	if err == nil && output.ErrorCode != 0 {
		err = fmt.Errorf("error executing API_SetDBvar: %s (error code: %v)", output.ErrorText, output.ErrorCode)
	}
	return
}

// UploadFileInput models the request sent to API_UploadFile
// See https://help.quickbase.com/api-guide/uploadfile.html
type UploadFileInput struct {
	RequestParams
	Credentials

	TableID  string                 `xml:"-"`
	Fields   []UploadFileInputField `xml:"field"`
	RecordID int                    `xml:"rid"`
}

func (input *UploadFileInput) setCredentials(creds Credentials) { input.Credentials = creds }
func (input *UploadFileInput) method() string                   { return http.MethodPost }
func (input *UploadFileInput) uri() string                      { return "/db/" + input.TableID }
func (input *UploadFileInput) payload() ([]byte, error)         { return xml.Marshal(input) }
func (input *UploadFileInput) headers(req *http.Request) {
	req.Header.Set("Content-Type", "application/xml")
	req.Header.Set("QUICKBASE-ACTION", "API_UploadFile")
}

// UploadFileInputField models the "field" element in API_UploadFile requests.
type UploadFileInputField struct {
	ID       int    `xml:"fid,attr"`
	FileData string `xml:",chardata"`
	Name     string `xml:"filename,attr"`
}

// UploadFileOutput models the response returned by API_UploadFile
// See https://help.quickbase.com/api-guide/uploadfile.html
type UploadFileOutput struct {
	ResponseParams

	Fields []UploadFileOutputField `xml:"file_fields>field" json:"fields"`
}

func (output *UploadFileOutput) parse(body []byte, res *http.Response) error {
	return parseXML(output, body, res)
}

// UploadFileOutputField models the "file_fields>field" element in
// API_UploadFile responses.
type UploadFileOutputField struct {
	ID  int    `xml:"id,attr" json:"field_id"`
	URL string `xml:"url" json:"url"`
}

// UploadFile makes an API_UploadFile call.
// See https://help.quickbase.com/api-guide/uploadfile.html
func (c Client) UploadFile(input *UploadFileInput) (output UploadFileOutput, err error) {
	err = c.Do(input, &output)
	if err == nil && output.ErrorCode != 0 {
		err = fmt.Errorf("error executing API_UploadFile: %s (error code: %v)", output.ErrorText, output.ErrorCode)
	}
	return
}