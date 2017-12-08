package aws

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type awsProvider struct {}

func AwsProvider () awsProvider {
	return awsProvider{}
}

func (p awsProvider) getEnvCreds () (string, string) {
	id := fmt.Sprintf("AWS_ACCESS_KEY_ID=%s", os.Getenv("ZCLOUD_AWS_KEY_ID"))
	secret := fmt.Sprintf("AWS_SECRET_ACCESS_KEY=%s", os.Getenv("ZCLOUD_AWS_SECRET_KEY"))
	return id, secret
}

func GetCmdArgs (cmd *exec.Cmd, args []string) []string {
	resArgs := cmd.Args
	for _, arg := range args {
		resArgs = append(resArgs, arg)
	}
	return resArgs
}

const cloudURLPrefix = "cloud://"

func ConvertURL (url string) string {
	if len(url) > len(cloudURLPrefix) && url[:len(cloudURLPrefix)] == cloudURLPrefix {
		return strings.Replace(url, "cloud", "s3", 1)
	}
	return url
}
