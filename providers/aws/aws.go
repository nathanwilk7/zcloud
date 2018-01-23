package aws

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

type awsProvider struct {}

func AwsProvider () awsProvider {
	return awsProvider{}
}

func getEnvCreds () (string, string, error) {
	id := os.Getenv("ZCLOUD_AWS_KEY_ID")
	secret := os.Getenv("ZCLOUD_AWS_SECRET_KEY")
	if id == "" || secret == "" {
		// TODO: err message
		return "", "", errors.New(fmt.Sprintf("env vars not found"))
	}
	return id, secret, nil
}

func getCreds () (string, string, error) {
	return getEnvCreds()
}

func getSession () *session.Session{
	id, secret, _ := getCreds()
	sess, _ := session.NewSession(
		&aws.Config{
			Region: aws.String("us-east-1"),
			Credentials: credentials.NewStaticCredentials(id, secret, ""),
		},
	)
	return sess
}

const cloudURLPrefix = "cloud://"

func convertURL (url string) string {
	if isCloudUrl(url) {
		return strings.Replace(url, "cloud", "s3", 1)
	}
	return url
}

func isCloudUrl (url string) bool {
	if len(url) > len(cloudURLPrefix) {
		if url[:len(cloudURLPrefix)] == cloudURLPrefix {
			return true
		}
	}
	return false
}

func awsStorageCmd (cmdStr string, urls []string, args []string) *exec.Cmd {
	keyId, secret, _ := getCreds()
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
