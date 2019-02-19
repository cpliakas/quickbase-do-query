package quickbase

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

// Plugin is the interface implemented by plugins that hook into the API
// request process.
type Plugin interface {

	// PreRequest is the hook that is invoked prior to the request being
	// sent to the Quick Base API.
	PreRequest(context.Context, *http.Request) context.Context

	// PostResponse is the hook invoked after the response is sent from the
	// Quick Base API. This hook is invoked whether the request was
	// successful or not.
	PostResponse(context.Context, *http.Request, *http.Response, []byte, error) context.Context
}

// Client makes requests to the Quick Base API.
type Client struct {

	// config stores the runtime configuration.
	Config Config

	// The HTTP client to use when sending requests. Defaults to
	// `http.DefaultClient`.
	HTTPClient *http.Client

	// Plugins contains the Plugin implementations.
	Plugins []Plugin
}

// NewClient returns a Client populated with default values.
func NewClient(cfg Config) Client {
	return Client{
		Config:     cfg,
		HTTPClient: http.DefaultClient,
	}
}

// NewRequest returns a *http.Request initialized with the data needed to
// make a request to the Quick Base API. In this method, the credentials are
// set, the Input struct is marshaled into the XML/JSON/HTML payload, and the
// URL of the action being performed is constructed.
func (c Client) NewRequest(input Input) (req *http.Request, err error) {

	if i, ok := input.(AuthenticatedInput); ok {
		i.setCredentials(NewCredentials(c.Config))
		input = i
	}

	b, err := input.payload()
	if err != nil {
		return
	}

	url := strings.TrimRight(c.Config.RealmHost(), "/") + input.uri()
	req, err = http.NewRequest(input.method(), url, bytes.NewBuffer(b))
	if err != nil {
		return
	}

	input.headers(req)
	return
}

// Do makes a request to the Quick Base API. This method sets some context
// about the request, initializes the request via the NewRequest method,
// invokes each plugins PreRequest method, uses *Client.HTTPClient to make
// the actual request, invokes each plugin's PostResponse method, then
// unmarshals the raw response into the passed Output struct.
func (c Client) Do(input Input, output Output) error {
	ctx := context.Background()
	ctx = context.WithValue(ctx, CtxKeyRealmHost, c.Config.RealmHost())

	req, err := c.NewRequest(input)
	if err != nil {
		return err
	}

	ctx = context.WithValue(ctx, CtxKeyAction, req.Header.Get("QUICKBASE-ACTION"))
	ctx = c.invokePreRequest(ctx, req)

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		ctx = c.invokePostResponse(ctx, req, res, []byte(""), err)
		return err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	ctx = c.invokePostResponse(ctx, req, res, body, err)
	if err != nil {
		return err
	}

	err = output.parse(body, res)
	if err != nil {
		return err
	}

	return nil
}

// invokePreRequest invokes each plugin's PreRequest method.
func (c Client) invokePreRequest(ctx context.Context, req *http.Request) context.Context {
	for _, p := range c.Plugins {
		ctx = p.PreRequest(ctx, req)
	}
	return ctx
}

// invokePostResponse invokes each plugin's PostResponse method.
func (c Client) invokePostResponse(ctx context.Context, req *http.Request, res *http.Response, body []byte, err error) context.Context {
	for _, p := range c.Plugins {
		ctx = p.PostResponse(ctx, req, res, body, err)
	}
	return ctx
}

// parseXML parses an XML response, populating output with data.
func parseXML(output Output, body []byte, res *http.Response) error {
	return xml.Unmarshal(body, output)
}

// parseJSON parses a JSON response, populating output with data.
func parseJSON(output Output, body []byte, res *http.Response) error {
	return json.Unmarshal(body, output)
}

// parseHTML parses an HTML response, populating output with data.
func parseHTML(output HTMLOutput, body []byte, res *http.Response) error {

	c, err := strconv.Atoi(res.Header.Get("QUICKBASE-ERRCODE"))
	if err != nil {
		return err
	}

	output.setErrorCode(c)
	output.setErrorText(res.Header.Get("QUICKBASE-ERRTEXT"))
	output.setHtml(body)

	return nil
}
