package cliutil

import "encoding/json"

// RenderJSON returns pretty-printed JSON as a string.
func RenderJSON(v interface{}) string {
	b, err := json.MarshalIndent(v, "", "    ")
	if err != nil {
		panic(err)
	}
	return string(b)
}
