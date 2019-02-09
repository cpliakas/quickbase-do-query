package cliutil

import (
	"fmt"
	"os"
)

// ErrorResponse models an error response.
type ErrorResponse struct {
	Error  string `json:"error"`
	Detail string `json:"detail,omitempty"`
}

// HandleError prints an error message in JSON format to STDERR and exits with
// a non-zero status code.
func HandleError(err error, prefix string) {
	if err == nil {
		return
	}

	resp := ErrorResponse{}
	if prefix != "" {
		resp.Error = prefix
		resp.Detail = err.Error()
	} else {
		resp.Error = err.Error()
	}

	fmt.Fprintf(os.Stderr, "%s\n", FormatJSON(resp))
	os.Exit(1)
}
