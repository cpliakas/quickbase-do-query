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

	// Split the string by spaces that aren't in quotes.
	parts := strings.FieldsFunc(s, fn)

	// Build and return the map.
	m := make(map[string]string)
	for _, part := range parts {
		p := strings.Split(part, "=")

		// Protect against values with no "=", treat them as a key.
		if len(p) < 2 {
			p = []string{p[0], ""}
		}

		// Trim quotes at the edges.
		// TODO parse to rune to check for unicode.Quotation_Mark?
		p[0] = strings.Trim(p[0], `"`)
		p[1] = strings.Trim(p[1], `"`)

		m[p[0]] = p[1]
	}

	return m
}
