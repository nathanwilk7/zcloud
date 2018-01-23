package aws

import (
	"errors"
	"fmt"
	"os"
	"strings"
	
	"github.com/nathanwilk7/zcloud/storage"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func (p awsProvider) Upload (params storage.CpParams) (string, error) {
	f, err  := os.Open(params.Src)
	if err != nil {
		return "", fmt.Errorf("failed to open file %s, %v", params.Src, err)
	}
	dest := convertURL(params.Dest)
	myBucket, myKey, err := getS3BucketKeyFromURL(dest)
	if err != nil {
		return "", err
	}
	sess, err := getSession()
	if err != nil {
		return "", err
	}
	uploader := s3manager.NewUploader(sess)
	result, err := uploader.Upload(
		&s3manager.UploadInput{
			Bucket: aws.String(myBucket),
			Key:    aws.String(myKey),
			Body:   f,
		},
	)
	if err != nil {
		return "", fmt.Errorf("failed to upload file, %v", err)
	}
	return fmt.Sprintf("%s uploaded to %s\n", params.Src, aws.StringValue(&result.Location)), nil
}

func (p awsProvider) Download (params storage.CpParams) (string, error) {
	return "", nil
}

func (p awsProvider) Cp (params storage.CpParams) (string, error) {
	if isCloudURL(params.Src) && !isCloudURL(params.Dest) {
		return p.Download(params)
	} else if !isCloudURL(params.Src) && isCloudURL(params.Dest) {
		return p.Upload(params)
	}
	return "", errors.New(fmt.Sprintf("Exactly one of the source and destination url's must be a cloud url: %s, %s", params.Src, params.Dest))
	//if params.Recursive {
	 //	panic("TODO recursive")
	//}
}

func (p awsProvider) Ls (params storage.LsParams) (string, error) {
	return "", nil
}

func (p awsProvider) Rm (params storage.RmParams) (string, error) {
	return "", nil
}

func (p awsProvider) Mv (params storage.MvParams) (string, error) {
	return "", nil
}

func (p awsProvider) Mb (params storage.MbParams) (string, error) {
	return "", nil
}

func (p awsProvider) Rb (params storage.RmParams) (string, error) {
	return "", nil
}

func (p awsProvider) Sync (params storage.SyncParams) (string, error) {
	return "", nil
}

const s3Prefix = "s3://"

func getS3BucketKeyFromURL (url string) (string, string, error) {
	if len(url) <= len(s3Prefix) {
		return "", "", fmt.Errorf("Converted URL %s must begin with %s followed by a bucket name", url, s3Prefix)
	}
	if !strings.Contains(url[len(s3Prefix):], "/") {
		return url[len(s3Prefix):], "", nil
	}
	firstSlash := strings.Index(url[len(s3Prefix):], "/") + len(s3Prefix)
	return url[len(s3Prefix):firstSlash], url[firstSlash + 1:], nil
}
