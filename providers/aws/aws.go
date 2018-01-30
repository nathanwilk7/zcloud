package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

type awsProvider struct {
	Id, Secret, Region string
	Session *session.Session
}

func AwsProvider (id, secret, region string) awsProvider {
	s, err := getSession(id, secret, region)
	if err != nil {
		panic(err)
	}
	return awsProvider{
		Id: id,
		Secret: secret,
		Region: region,
		Session: s,
	}
}

const defaultToken = ""

func getSession (id, secret, region string) (*session.Session, error) {
	sess, err := session.NewSessionWithOptions(
		session.Options{
			Config: aws.Config{
				Region: aws.String(region),
				Credentials: credentials.NewStaticCredentials(id, secret, defaultToken),
			},
		},
	)
	if err != nil {
		return nil, err
	}
	return sess, nil
}
