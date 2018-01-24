package providers

import (
	"fmt"
	"os"
	
	"github.com/nathanwilk7/zcloud/storage"
	"github.com/nathanwilk7/zcloud/compute"

	"github.com/nathanwilk7/zcloud/providers/aws"
	"github.com/nathanwilk7/zcloud/providers/gcloud"
	"github.com/nathanwilk7/zcloud/providers/test_provider"
)

type Provider interface {
	storage.StorageProvider
	compute.ComputeProvider
}

const (
	ZCloudProvEnv = "ZCLOUD_PROV"
	ZCloudStorageProvEnv = "ZCLOUD_STORAGE_PROV"
)

func GetProvider () (Provider, error) {
	return getProvider()
}

func getProvider() (Provider, error) {
	prov := os.Getenv(ZCloudProvEnv)
	if p, ok := providers[prov]; ok {
		return p, nil
	}
	return nil, fmt.Errorf("%s was not valid or was empty: %s", ZCloudProvEnv, prov)
}

var providers map[string]Provider = map[string]Provider {
	"TEST": test_provider.TestProvider(),
	"GCLOUD": gcloud.GCloudProvider(),
	"AWS": aws.AwsProvider(),
}

func GetStorageProvider () (storage.StorageProvider, error) {
	return getStorageProvider()
}

func getStorageProvider() (storage.StorageProvider, error) {
	prov := os.Getenv(ZCloudStorageProvEnv)
	if p, ok := storageProviders[prov]; ok {
		return p, nil
	}
	if p, ok := providers[prov]; ok {
		return p, nil
	}
	if p, err := getProvider(); err == nil {
		return p, nil
	}
	return nil, fmt.Errorf("%s and %s were not valid or were empty: %s, %s", ZCloudProvEnv, ZCloudStorageProvEnv, prov, os.Getenv(ZCloudProvEnv))
}

var storageProviders map[string]storage.StorageProvider = map[string]storage.StorageProvider {}
