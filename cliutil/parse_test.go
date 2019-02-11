package cliutil_test

import (
	"testing"

	"github.com/cpliakas/quickbase-do-query/cliutil"
)

func TestParseKeyValue(t *testing.T) {
	s := `time="2017-05-30T19:02:08-05:00" level=info msg="some log message"`
	m := cliutil.ParseKeyValue(s)

	tests := []struct {
		key  string
		want string
	}{
		{"time", "2017-05-30T19:02:08-05:00"},
		{"level", "info"},
		{"msg", "some log message"},
	}

	for _, test := range tests {
		if m[test.key] != test.want {
			t.Errorf("got '%s', want '%s'", m[test.key], test.want)
		}
	}
}
