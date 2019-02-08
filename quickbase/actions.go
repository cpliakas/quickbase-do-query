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

	TableID          string                `xml:"-"`
	Query            string                `xml:"query,omitempty"`
	QueryID          int                   `xml:"qid,omitempty"`
	QueryName        string                `xml:"qname,omitempty"`
	IncludeRecordIDs BoolToInt             `xml:"includeRids,omitempty"`
	ReturnPercentage BoolToInt             `xml:"returnpercentage,omitempty"`
	UseFIDs          BoolToInt             `xml:"useFids,omitempty"`
	Format           string                `xml:"fmt,omitempty"`
	FieldSlice       DoQueryInput_Fields   `xml:"clist,omitempty"`
	SortSlice        DoQueryInput_Fields   `xml:"slist,omitempty"`
	Options          *DoQueryInput_Options `xml:"options,omitempty"`
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
func (i *DoQueryInput) EnsureOptions() *DoQueryInput_Options {
	if i.Options == nil {
		i.Options = &DoQueryInput_Options{}
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

// DoQueryInput_Fields models field lists in API_DoQuery requests.
type DoQueryInput_Fields []int

// MarshalXML converts a list of fields to a "." delimited string.
func (f DoQueryInput_Fields) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.EncodeElement(FormatFieldIDs(f), start)
}

// DoQueryInput_Options models the "options" element in API_DoQuery requests.
type DoQueryInput_Options struct {
	SortOrderSlice []string
	Limit          int
	Offset         int
	OnlyNew        bool
	Unsorted       bool
}

// MarshalXML implements Marshaler.MarshalXML and formats the value of the
// "options" element.
func (o DoQueryInput_Options) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
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

	Fields  []DoQueryOutput_Field  `xml:"table>fields>field"`
	Records []DoQueryOutput_Record `xml:"table>records>record"`
}

// DoQueryOutput_Field models the "table>fields" element. in an API_DoQuery
// response.
type DoQueryOutput_Field struct {
	FieldID  int    `xml:"id,attr"`
	Type     string `xml:"field_type,attr"`
	BaseType string `xml:"base_type,attr"`
	Mode     string `xml:"mode,attr"`
	Label    string `xml:"label"`
}

// DoQueryOutput_Record models the "table>records>record" element in an
// API_DoQuery response.
type DoQueryOutput_Record struct {
	RecordID int                          `xml:"rid,attr"`
	UpdateID int                          `xml:"update_id"`
	Fields   []DoQueryOutput_Record_Field `xml:"f"`
}

// DoQueryOutput_Record_Field models the "table>records>record>field" element
// in an API_DoQuery response.
type DoQueryOutput_Record_Field struct {
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

	Fields []DoQueryOutput_Field `xml:"table>fields>field"`
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
