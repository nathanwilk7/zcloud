package test_provider

import (
	"fmt"
	"os"
	
	"github.com/nathanwilk7/zcloud/storage"
)

func (p testProvider) Upload (params storage.UploadParams) (string, error) {
	if _, err := os.Stat(params.Src); os.IsNotExist(err) {
		return "", fmt.Errorf("Source file does not exist: %s", params.Src)
	}
	return "", nil
}

func (p testProvider) Download (params storage.DownloadParams) (string, error) {
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

const testStr = "test"
const testPrefix = testStr + "://"

func (p testProvider) StorageURLPrefixReplacement() string {
	return testStr
}
