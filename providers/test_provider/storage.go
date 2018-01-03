package test_provider

import (
	"github.com/nathanwilk7/zcloud/storage"
)

func (p testProvider) Cp (params storage.CpParams) (string, error) {
	return "", nil
}

func (p testProvider) Ls (params storage.LsParams) (string, error) {
	return "", nil
}
