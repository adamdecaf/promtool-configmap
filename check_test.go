package main

import (
	"os"
	"testing"
)

func TestCheck(t *testing.T) {
	cases := []struct {
		path  string
		valid bool
	}{
		{
			path:  "testdata/bad.yaml",
			valid: false,
		},
		{
			path:  "testdata/bad.json",
			valid: false,
		},
		{
			path:  "testdata/empty.yaml",
			valid: false,
		},
		{
			path:  "testdata/good.yaml",
			valid: true,
		},
		{
			path:  "testdata/good.json",
			valid: true,
		},
	}
	for i := range cases {
		f, err := os.Open(cases[i].path)
		if err != nil {
			t.Fatalf("error reading %s, err=%v", cases[i].path, err)
		}

		err = check(f)
		if cases[i].valid && err != nil {
			t.Errorf("got error on %s, err=%v", cases[i].path, err)
		}
		if !cases[i].valid && err == nil {
			t.Errorf("expected error on %s", cases[i].path)
		}
	}
}
