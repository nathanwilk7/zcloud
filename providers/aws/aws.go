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



func getEnvCreds () (string, string, error) {
	idEnv := os.Getenv("ZCLOUD_AWS_KEY_ID")
	secretEnv := os.Getenv("ZCLOUD_AWS_SECRET_KEY")
	if idEnv == nil || secretEnv == nil {
		return "", "", errors.New(fmt.Sprintf("
}
	id := fmt.Sprintf("AWS_ACCESS_KEY_ID=%s", idEnv)
	secret := fmt.Sprintf("AWS_SECRET_ACCESS_KEY=%s", secretEnv)
	return id, secret, nil
}

// func getCmdArgs (cmd *exec.Cmd, args []string) []string {
// 	resArgs := cmd.Args
// 	for _, arg := range args {
// 		resArgs = append(resArgs, arg)
// 	}
// 	return resArgs
// }

const cloudURLPrefix = "cloud://"

func convertURL (url string) string {
	if len(url) > len(cloudURLPrefix) && url[:len(cloudURLPrefix)] == cloudURLPrefix {
		return strings.Replace(url, "cloud", "s3", 1)
	}
	return url
}

func awsStorageCmd (cmdStr string, urls []string, args []string) *exec.Cmd {
	keyId, secret := getEnvCreds()
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
