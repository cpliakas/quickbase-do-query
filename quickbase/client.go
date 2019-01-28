package quickbase

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// Client makes requests to the Quick Base API.
type Client struct {
	Config  *Config
	TableID string
}

// DoQuery executes an API_DoQuery request.
func (c *Client) DoQuery(in DoQueryInput) (DoQueryOutput, error) {
	out := DoQueryOutput{}

	// Add credentials.
	if c.Config.UserToken != "" {
		in.UserToken = c.Config.UserToken
	} else if c.Config.Ticket != "" {
		in.Ticket = c.Config.Ticket
		if c.Config.AppToken != "" {
			in.AppToken = c.Config.AppToken
		}
	}

	// Set required parameters.
	in.Format = "structured"
	in.IncludeRecordIDs = true
	in.UseFIDs = false

	// Format the XML payload.
	payload, err := xml.Marshal(in)
	if err != nil {
		return out, err
	}

	// Build the HTTP request, add required headers.
	url := fmt.Sprintf("%s/db/%s", strings.TrimRight(c.Config.RealmHost, "/"), c.TableID)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(payload))
	if err != nil {
		return out, err
	}

	req.Header.Add("QUICKBASE-ACTION", "API_DoQuery")
	req.Header.Add("Content-Type", "application/xml")

	// Create an http.Client if one isn't set.
	if c.Config.HTTPClient == nil {
		c.Config.HTTPClient = &http.Client{
			Timeout: 5 * time.Second,
		}
	}

	// Execute the API request.
	resp, err := c.Config.HTTPClient.Do(req)
	if err != nil {
		return out, err
	}

	// Parse the response.
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return out, err
	}

	err = xml.Unmarshal(body, &out)
	return out, err
}
