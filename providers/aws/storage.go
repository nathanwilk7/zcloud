package aws

import (
	"errors"
	"fmt"
	"os"
	"strings"
	
	"github.com/nathanwilk7/zcloud/storage"

	//	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func (p awsProvider) Upload (params storage.CpParams) (string, error) {	
	f, err  := os.Open(params.Src)
	if err != nil {
		return "", fmt.Errorf("failed to open file %q, %v", params.Src, err)
	}
	dest := convertURL(params.Dest)
	myBucket, myKey, _ := getS3BucketKeyFromUrl(dest)
	uploader := s3manager.NewUploader(getSession())
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
	fmt.Printf("file uploaded to, %s\n", aws.StringValue(&result.Location))
	
	return "", nil
}

func (p awsProvider) Download (params storage.CpParams) (string, error) {
	return "", nil
}

func (p awsProvider) Cp (params storage.CpParams) (string, error) {
	if isCloudUrl(params.Src) && !isCloudUrl(params.Dest) {
		return p.Download(params)
	} else if !isCloudUrl(params.Src) && isCloudUrl(params.Dest) {
		return p.Upload(params)
	}
	return "", errors.New(fmt.Sprintf("Exactly one of the source and destination url's must be a cloud url"))
	// if params.Recursive {
	// 	panic("TODO recursive")
	// }
}

func (p awsProvider) Ls (params storage.LsParams) (string, error) {
	args := []string{}
	if params.Recursive {
		args = append(args, "--recursive")
	}
	cmd := awsStorageCmd(
		"ls",
		[]string{params.Url},
		args,
	)
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return "List completed successfully", nil
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

func getS3BucketKeyFromUrl (url string) (string, string, error) {
	if len(url) <= len(s3Prefix) {
		return "", "", fmt.Errorf("TODO error")
	}
	if !strings.Contains(url[len(s3Prefix):], "/") {
		return url[len(s3Prefix):], "", nil
	}
	firstSlash := strings.Index(url[len(s3Prefix):], "/") + len(s3Prefix)
	return url[len(s3Prefix):firstSlash], url[firstSlash + 1:], nil
}
