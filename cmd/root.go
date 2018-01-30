package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/nathanwilk7/zcloud/providers"
	"github.com/nathanwilk7/zcloud/providers/aws"
	"github.com/nathanwilk7/zcloud/providers/gcloud"
	"github.com/nathanwilk7/zcloud/providers/test_provider"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "zcloud",
	Short: "zCloud makes using the cloud better",
	Long: "zCloud provides a layer of abstraction between clients and cloud providers",
}

var quiet bool

func init () {
	RootCmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false, "Don't print output")
}

func writeOutput (msg string) {
	if !quiet {
		fmt.Println(msg)
	}
}

const (
	provEnv = "ZCLOUD_PROV"
	awsIdEnv = "ZCLOUD_AWS_KEY_ID"
	awsSecretEnv = "ZCLOUD_AWS_SECRET_KEY"
	awsRegionEnv = "ZCLOUD_AWS_REGION"
)

func getProvider (provider string) (providers.Provider, error) {
	if p, ok := providersMap[provider]; ok {
		return p, nil
	}
	return nil, fmt.Errorf("%s was not a valid Provider", provider)
}

// TODO: Break out initialization of providers so that we only initialize the provider we're using. Also, do error checking.
var providersMap map[string]providers.Provider = map[string]providers.Provider {
	"TEST": test_provider.TestProvider(),
	"GCLOUD": gcloud.GCloudProvider(),
	"AWS": aws.AwsProvider(
		os.Getenv(awsIdEnv),
		os.Getenv(awsSecretEnv),
		os.Getenv(awsRegionEnv),
	),
}

const (
	cloudStr = "cloud"
	cloudURLPrefix = cloudStr + "://"
)

func convertURL (url, replacement string) string {
	if isCloudURL(url) {
		return strings.Replace(url, cloudStr, replacement, 1)
	}
	return url
}

func isCloudURL (url string) bool {
	if len(url) > len(cloudURLPrefix) {
		if url[:len(cloudURLPrefix)] == cloudURLPrefix {
			return true
		}
	}
	return false
}
