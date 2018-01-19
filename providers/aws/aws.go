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

func getCreds () (string, string) {
	return getEnvCreds()
}

func getEnvCreds () (string, string) {
	id := fmt.Sprintf("AWS_ACCESS_KEY_ID=%s", os.Getenv("ZCLOUD_AWS_KEY_ID"))
	secret := fmt.Sprintf("AWS_SECRET_ACCESS_KEY=%s", os.Getenv("ZCLOUD_AWS_SECRET_KEY"))
	return id, secret
}

const cloudURLPrefix = "cloud://"

func convertURL (url string) string {
	if len(url) > len(cloudURLPrefix) && url[:len(cloudURLPrefix)] == cloudURLPrefix {
		return strings.Replace(url, "cloud", "s3", 1)
	}
	return url
}

func awsStorageCmd (cmdStr string, urls []string, args []string) *exec.Cmd {
	keyId, secret := getCreds()
	cmd := exec.Command("aws")
	cmd.Env = []string{keyId, secret}
	cmdArgs := []string{"s3", cmdStr}
	for _, arg := range args {
		cmdArgs = append(cmdArgs, arg)
	}
	for _, url := range urls {
		cmdArgs = append(cmdArgs, convertURL(url))
	}
	cmd.Args = cmdArgs
	return cmd
}
