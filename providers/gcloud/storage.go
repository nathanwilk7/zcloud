package gcloud

import (
	"github.com/nathanwilk7/zcloud/storage"
)

func (p gcloudProvider) Cp (params storage.CpParams) (string, error) {
	return "", nil
}

func (p gcloudProvider) Ls (params storage.LsParams) (string, error) {
	return "", nil
}

func (p gcloudProvider) Rm (params storage.RmParams) (string, error) {
	return "", nil
}

func (p gcloudProvider) Mv (params storage.MvParams) (string, error) {
	return "", nil
}

func (p gcloudProvider) Mb (params storage.MbParams) (string, error) {
	return "", nil
}

func (p gcloudProvider) Rb (params storage.RmParams) (string, error) {
	return "", nil
}

func (p gcloudProvider) Sync (params storage.SyncParams) (string, error) {
	return "", nil
}
