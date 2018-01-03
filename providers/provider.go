package providers

import (
	"errors"
	"fmt"

	"os"
	
	"github.com/nathanwilk7/zcloud/storage"
	"github.com/nathanwilk7/zcloud/compute"

	"github.com/nathanwilk7/zcloud/providers/test_provider"
	"github.com/nathanwilk7/zcloud/providers/aws"
	"github.com/nathanwilk7/zcloud/providers/gcloud"
)

type Provider interface {
	storage.StorageProvider
	compute.ComputeProvider
}

const (
	ZCloudProvEnv = "ZCLOUD_PROV"
	ZCloudStorageProvEnv = "ZCLOUD_STORAGE_PROV"
)

func getProvEnv() string {
	return os.Getenv(ZCloudProvEnv)
}

// TODO: How to avoid duplication between getProvider and getStorageProvider?
func getProvider() (Provider, error) {
	prov := getProvEnv()
	switch prov {
	case "TEST":
		return test_provider.TestProvider(), nil
	// case "AWS":
	// 	return aws.AwsProvider(), nil
	case "GCLOUD":
		return gcloud.GCloudProvider(), nil
	default:
		return nil, errors.New(fmt.Sprintf("%s was not valid or was empty: %s", ZCloudProvEnv, prov))
	}
}

func GetProvider() (Provider, error) {
	return getProvider()
}

func getStorageProvider() (storage.StorageProvider, error) {
	prov := os.Getenv(ZCloudStorageProvEnv)
	if prov == "" {
		prov = getProvEnv()
	}
	switch prov {
	case "TEST":
		return test_provider.TestProvider(), nil
	case "AWS":
		return aws.AwsProvider(), nil
	case "GCLOUD":
		return gcloud.GCloudProvider(), nil
	default:
		return nil, errors.New(fmt.Sprintf("%s and %s were not valid or were empty: %s", ZCloudProvEnv, ZCloudStorageProvEnv, prov))
	}
}

func GetStorageProvider() (storage.StorageProvider, error) {
	return getStorageProvider()
}
