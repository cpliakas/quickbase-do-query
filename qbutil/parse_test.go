package qbutil_test

import (
	"testing"

	"github.com/cpliakas/quickbase-do-query/qbutil"
)

func TestSplit(t *testing.T) {
	s := "1,2.3 , 4 ASC.  5"
	parts := qbutil.Split(s)

	if len(parts) != 5 {
		t.Fatalf("expected 5 fields, got %v", len(parts))
	}
	if parts[0] != "1" {
		t.Errorf("expected '1', got '%s'", parts[0])
	}
	if parts[1] != "2" {
		t.Errorf("expected '2', got '%s'", parts[1])
	}
	if parts[2] != "3" {
		t.Errorf("expected '3', got '%s'", parts[2])
	}
	if parts[3] != "4 ASC" {
		t.Errorf("expected '4 ASC', got '%s'", parts[3])
	}
	if parts[4] != "5" {
		t.Errorf("expected '5', got '%s'", parts[4])
	}
}
