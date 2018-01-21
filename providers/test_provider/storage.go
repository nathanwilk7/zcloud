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

func (p testProvider) Rm (params storage.RmParams) (string, error) {
	return "", nil
}

func (p testProvider) Mv (params storage.MvParams) (string, error) {
	return "", nil
}

func (p testProvider) Mb (params storage.MbParams) (string, error) {
	return "", nil
}

func (p testProvider) Rb (params storage.RmParams) (string, error) {
	return "", nil
}

func (p testProvider) Sync (params storage.SyncParams) (string, error) {
	return "", nil
}
