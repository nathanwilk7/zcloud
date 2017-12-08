package providers

import (
	"errors"
	"fmt"
	"log"
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

const ZCloudProvEnv = "ZCLOUD_PROV" 

func getProvider() (Provider, error) {
	prov := os.Getenv(ZCloudProvEnv)
	switch prov {
	case "TEST":
		return test_provider.TestProvider(), nil
	case "AWS":
		return aws.AwsProvider(), nil
	case "GCLOUD":
		return gcloud.GCloudProvider(), nil
	default:
		return nil, errors.New(fmt.Sprintf("%s was not valid or was empty: %s", ZCloudProvEnv, prov))
	}
}

var ProviderInstance Provider

func init () {
	providerInstance, err := getProvider()
	if err != nil {
		log.Fatal(err)
	}
	ProviderInstance = providerInstance
}
