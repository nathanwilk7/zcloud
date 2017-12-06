package main

import (
	"github.com/nathanwilk7/zcloud/test_provider"
	// "github.com/nathanwilk7/zcloud/aws"
	// "github.com/nathanwilk7/zcloud/gcloud"
)

type Provider interface {
	StorageProvider
	ComputeProvider
}

func getProvider() Provider {
	return test_provider.TestProvider{}
}

func main () {
	// parse command, https://blog.komand.com/build-a-simple-cli-tool-with-golang
	// create provider
	p := getProvider() // env vars specify creds
	// call appropriate provider func based on parsed command
	p.Cp("test.txt", "s3://cloud-agnostic-testing/test.txt")
}
