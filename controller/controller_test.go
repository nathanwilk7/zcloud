package controller

import (
	"testing"
)

func TestBaseNameOrDir (t *testing.T) {
	f := "a.txt"
	k := ""
	s := baseNameOrDir(f, k)
	if s != "a.txt" {
		t.Fatal(f, k, s)
	}
}
	]
