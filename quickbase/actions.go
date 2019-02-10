package quickbase

import (
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
	IncludeRecordIDs BoolToInt            `xml:"includeRids,omitempty"`
	ReturnPercentage BoolToInt            `xml:"returnpercentage,omitempty"`
	UseFIDs          BoolToInt            `xml:"useFids,omitempty"`
	Format           string               `xml:"fmt,omitempty"`
	FieldSlice       DoQueryInputFields   `xml:"clist,omitempty"`
	SortSlice        DoQueryInputFields   `xml:"slist,omitempty"`
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
func (i *DoQueryInput) EnsureOptions() *DoQueryInputOptions {
	if i.Options == nil {
		i.Options = &DoQueryInputOptions{}
	}
	return i.Options
}

// Fields sets the fields that are returned in the response.
func (i *DoQueryInput) Fields(fids ...int) *DoQueryInput {
	i.FieldSlice = fids
	return i
}

// SortBy sets the fields to be sorted by.
func (i *DoQueryInput) SortBy(fids ...int) *DoQueryInput {
	i.SortSlice = fids
	return i
}

// SortOrder sets the "sortorder" option.
func (i *DoQueryInput) SortOrder(order ...string) *DoQueryInput {
	i.EnsureOptions().SortOrderSlice = order
	return i
}

// Sort sets the fields and order in one shot.
func (i *DoQueryInput) Sort(sort []int, order []string) *DoQueryInput {
	i.SortSlice = sort
	i.EnsureOptions().SortOrderSlice = order
	return i
}

// Limit sets the "num" option.
func (i *DoQueryInput) Limit(n int) *DoQueryInput {
	i.EnsureOptions().Limit = n
	return i
}

// Offset sets the "skp" option.
func (i *DoQueryInput) Offset(n int) *DoQueryInput {
	i.EnsureOptions().Offset = n
	return i
}

// OnlyNew sets the "onlynew" option.
func (i *DoQueryInput) OnlyNew() *DoQueryInput {
	i.EnsureOptions().OnlyNew = true
	return i
}

// Unsorted sets the "nosort" option.
func (i *DoQueryInput) Unsorted() *DoQueryInput {
	i.EnsureOptions().Unsorted = true
	return i
}

// DoQueryInputFields models field lists in API_DoQuery requests.
type DoQueryInputFields []int

// MarshalXML converts a list of fields to a "." delimited string.
func (f DoQueryInputFields) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.EncodeElement(FormatFieldIDs(f), start)
}

// DoQueryInputOptions models the "options" element in API_DoQuery requests.
type DoQueryInputOptions struct {
	SortOrderSlice []string
	Limit          int
	Offset         int
	OnlyNew        bool
	Unsorted       bool
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
	if len(o.SortOrderSlice) > 0 {
		opts = append(opts, "sortorder-"+strings.Join(o.SortOrderSlice, ""))
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

// UploadFileInputField models the "field" element in an
// API_UploadFile request.
type UploadFileInputField struct {
	FieldID  int    `xml:"fid,attr"`
	FileData string `xml:",chardata"`
	Name     string `xml:"filename,attr"`
}

// UploadFileOutput models the response returned by API_UploadFile
// See https://help.quickbase.com/api-guide/uploadfile.html
type UploadFileOutput struct {
	ResponseParams

	Fields []UploadFileOutputField `xml:"file_fields>field"`
}

func (output *UploadFileOutput) parse(body []byte, res *http.Response) error {
	return parseXML(output, body, res)
}

// UploadFileOutputField models the "file_fields>field" element an
// API_UploadFile response.
type UploadFileOutputField struct {
	ID  int    `xml:"id,attr"`
	URL string `xml:"url"`
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
