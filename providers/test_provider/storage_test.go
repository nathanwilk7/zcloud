package test_provider

import (
	"testing"

	"github.com/nathanwilk7/zcloud/storage"
)

func TestCp (t *testing.T) {
	p := testProvider{}
	params := storage.CpParams{}
	if _, err := p.Cp(params); err != nil {
		t.Fatal(err)
	}
}
