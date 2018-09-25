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
		// empty config, never valid
		{
			path:  "testdata/empty.yml",
			valid: false,
		},
		// rule files
		{
			path:  "testdata/rules-bad.yml",
			valid: false,
		},
		{
			path:  "testdata/rules-bad.json",
			valid: false,
		},
		{
			path:  "testdata/rules-good.yml",
			valid: true,
		},
		{
			path:  "testdata/rules-good.json",
			valid: true,
		},
		// prom configs
		{
			path:  "testdata/config-good.yml",
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
