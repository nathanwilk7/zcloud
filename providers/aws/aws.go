package aws

import (
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

type awsProvider struct {}

func AwsProvider () awsProvider {
	return awsProvider{}
}

const idEnv = "ZCLOUD_AWS_KEY_ID"
const secretEnv = "ZCLOUD_AWS_KEY_ID"

func getCreds () (string, string, error) {
	return getEnvCreds()
}

func getEnvCreds () (string, string, error) {
	id := os.Getenv(idEnv)
	secret := os.Getenv(secretEnv)
	if id == "" || secret == "" {
		return "", "", fmt.Errorf("%s or %s was empty", idEnv, secretEnv)
	}
	return id, secret, nil
}

const regionEnv = "ZCLOUD_AWS_REGION"
const defaultToken = ""

func getSession () (*session.Session, error) {
	id, secret, err := getCreds()
	if err != nil {
		return nil, err
	}
	region := os.Getenv(regionEnv)
	if region == "" {
		return nil, fmt.Errorf("%s was empty", regionEnv)
	}
	sess, err := session.NewSession(
		&aws.Config{
			Region: aws.String(region),
			Credentials: credentials.NewStaticCredentials(id, secret, defaultToken),
		},
	)
	if err != nil {
		return nil, err
	}
	return sess, nil
}

const cloudStr = "cloud"
const cloudURLPrefix = cloudStr + "://"

func convertURL (url string) string {
	if isCloudURL(url) {
		return strings.Replace(url, cloudStr, "s3", 1)
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
