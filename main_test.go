package main

import "testing"

func TestScanFile(t *testing.T) {
	tests := []struct {
		filename string
		comments bool
		ok       bool
	}{
		// no homoglyphs
		{"testdata/test1.go", false, true},
		{"testdata/test1.go", true, true},
		// homoglyphs in source
		{"testdata/test2.go", false, false},
		{"testdata/test2.go", true, false},
		// homoglyphs in comments
		{"testdata/test3.go", false, true},
		{"testdata/test3.go", true, false},
	}

	for _, test := range tests {
		ok := len(scanFile(test.filename, test.comments)) == 0
		if ok != test.ok {
			t.Errorf("%v: expected %v, got %v", test.filename, test.ok, ok)
		}
	}
}
