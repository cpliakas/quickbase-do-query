package cliutil

import (
	"strings"
	"unicode"
)

// ParseKeyValue parses key value pairs.
// See https://stackoverflow.com/a/44282136
func ParseKeyValue(s string) map[string]string {

	lastQuote := rune(0)
	fn := func(c rune) bool {
		switch {
		case c == lastQuote:
			lastQuote = rune(0)
			return false
		case lastQuote != rune(0):
			return false
		case unicode.In(c, unicode.Quotation_Mark):
			lastQuote = c
			return false
		default:
			return unicode.IsSpace(c)
		}
	}

	// Split string by spaces that aren't in quotes.
	parts := strings.FieldsFunc(s, fn)

	// Build and return the map.
	m := make(map[string]string)
	for _, part := range parts {
		p := strings.Split(part, "=")
		m[p[0]] = strings.Trim(p[1], `"`) // TODO unicode.Quotation_Mark?
	}

	return m
}
