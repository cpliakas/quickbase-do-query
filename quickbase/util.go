package quickbase

import (
	"encoding/xml"
	"fmt"
	"strings"
)

// BoolToInt converts a boolean to an integer.
type BoolToInt bool

// MarshalXML implements Marshaler.MarshalXML and renders the bool as an int.
func (b BoolToInt) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	var s string
	if b == true {
		s = "1"
	} else {
		s = "0"
	}
	return e.EncodeElement(s, start)
}

// sliceToString convers a slice of integers into a "." delimited string.
func sliceToString(v []int) string {
	return strings.Trim(strings.Replace(fmt.Sprint(v), " ", ".", -1), "[]")
}