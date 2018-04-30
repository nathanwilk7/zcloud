package cmd

import (
	"os"
	
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

const (
	provEnv = "ZCLOUD_PROV"
	awsIdEnv = "ZCLOUD_AWS_KEY_ID"
	awsSecretEnv = "ZCLOUD_AWS_SECRET_KEY"
	awsRegionEnv = "ZCLOUD_AWS_REGION"
	gCloudProjectIDEnv = "ZCLOUD_GCLOUD_PROJECT_ID"
	destProvEnv = "ZCLOUD_DEST_PROV"
)

var (
	prov = os.Getenv(provEnv)
	awsId = os.Getenv(awsIdEnv)
	awsSecret = os.Getenv(awsSecretEnv)
	awsRegion = os.Getenv(awsRegionEnv)
	gCloudProjectID = os.Getenv(gCloudProjectIDEnv)
	destProv = os.Getenv(destProvEnv)
)
