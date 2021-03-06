package cliutil

import (
	"encoding/json"
	"fmt"
)

// FormatJSON returns pretty-printed JSON as a string.
func FormatJSON(v interface{}) string {
	b, err := json.MarshalIndent(v, "", "    ")
	if err != nil {
		panic(err)
	}
	return string(b)
}

// PrintJSON writes pretty-printed JSON to STDOUT.
func PrintJSON(v interface{}) {
	fmt.Println(FormatJSON(v))
}
