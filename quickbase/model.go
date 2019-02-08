package quickbase

import (
	"encoding/xml"
	"net/http"
)

// Input is the interface implemented by structs that model requests sent to
// the Quick Base API.
type Input interface {

	// method is HTTP method used for the request.
	method() string

	// uri is the path of the request, e.g. "/db/main"
	uri() string

	// headers adds headers to the HTTP request.
	headers(req *http.Request)

	// payload is the payload send as the body of the request.
	payload() ([]byte, error)
}

// RequestParams models the parameters that are common to all API requests.
type RequestParams struct {
	XMLName  xml.Name `xml:"qdbapi"`
	UserData string   `xml:"udata,omitempty"`
}

// AuthenticatedInput is the interface implemented by structs that model
// requests sent to the Quick Base API that require authentication.
type AuthenticatedInput interface {
	Input

	// setCredentials stores Credentials used to authenticate the request.
	setCredentials(Credentials)
}

// Output is the interface implemented by structs that model responses
// returned from Quick Base API requests.
type Output interface {

	// parse unmarshals the response body in the struct.
	parse(body []byte, resp *http.Response) error

	// setAction sets the action of the API request.
	setAction(string)

	// setErrorCode sets the numeric error code returned by Quick Base.
	setErrorCode(int)

	// setErrorText sets the error message returned by Quick Base.
	setErrorText(string)

	// setErrorText sets the detailed error message returned by Quick Base.
	setErrorDetail(string)
}

// HTMLOutput is the interfaces implemented by structs that model responses
// returned from the Quick Base API requests with raw HTML payloads.
type HTMLOutput interface {
	Output

	// setHtml sets the resonse body as a property in the struct.
	setHtml([]byte)
}

// ResponseParams implements Output and models the parameters that are
// common to all responses.
type ResponseParams struct {
	XMLName     xml.Name `xml:"qdbapi" json:"-"`
	Action      string   `xml:"action" json:"-"`
	ErrorCode   int      `xml:"errcode" json:"-"`
	ErrorText   string   `xml:"errtext" json:"-"`
	ErrorDetail string   `xml:"errdetail" json:"-"`
	UserData    string   `xml:"udata,omitempty" json:",omitempty"`
}

func (r *ResponseParams) setAction(a string)      { r.Action = a }
func (r *ResponseParams) setErrorCode(c int)      { r.ErrorCode = c }
func (r *ResponseParams) setErrorText(t string)   { r.ErrorText = t }
func (r *ResponseParams) setErrorDetail(d string) { r.ErrorDetail = d }
