package aws

import (
	"fmt"
	"os"
	"strings"

	"github.com/nathanwilk7/zcloud/storage"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func (p awsProvider) Upload (params storage.UploadParams) (string, error) {
	f, err  := os.Open(params.Src)
	if err != nil {
		return "", fmt.Errorf("failed to open file %s, %v", params.Src, err)
	}
	defer f.Close()
	bucket, key, err := bucketKeyFromURL(params.Dest)
	if err != nil {
		return "", err
	}
	uploader := s3manager.NewUploader(p.Session)
	result, err := uploader.Upload(
		&s3manager.UploadInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
			Body:   f,
		},
	)
	if err != nil {
		return "", fmt.Errorf("failed to upload file, %v", err)
	}
	return fmt.Sprintf("%s uploaded to %s", params.Src, aws.StringValue(&result.Location)), nil
}

func (p awsProvider) Download (params storage.DownloadParams) (string, error) {
	downloader := s3manager.NewDownloader(p.Session)
	
	f, err := os.Create(params.Dest)
	if err != nil {
		return "", fmt.Errorf("failed to create file %q, %v", params.Dest, err)
	}
	defer f.Close()

	bucket, key, err := bucketKeyFromURL(params.Src)
	if err != nil {
		return "", err
	}
	_, err = downloader.Download(f, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return "", fmt.Errorf("failed to download file %v", err)
	}
	return fmt.Sprintf("download: %s to %s", params.Src, params.Dest), nil
}

func (p awsProvider) Ls (params storage.LsParams) (string, error) {
	svc := s3.New(p.Session)
	bucket, key, err := bucketKeyFromURL(params.Url)
	if err != nil {
		return "", err
	}
	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
		Prefix: aws.String(key),
	}
	result, err := svc.ListObjectsV2(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeNoSuchBucket:
				return "", fmt.Errorf(s3.ErrCodeNoSuchBucket, aerr.Error())
			default:
				return "", aerr
			}
		} else {
			return "", err
		}
		return "", err
	}
	out := []string{}
	for _, content := range result.Contents {
		key := *content.Key
		if strings.Contains(key, "/") {
			out = append(out, key[:strings.Index(key, "/") + 1])
		} else {
			out = append(out, key)
		}
	}
	return strings.Join(out, "\n"), nil
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

const s3Str = "s3"
const s3Prefix = s3Str + "://"

func (p awsProvider) StorageURLPrefixReplacement() string {
	return s3Str
}

func bucketKeyFromURL (url string) (string, string, error) {
	if len(url) <= len(s3Prefix) {
		return "", "", fmt.Errorf("Converted URL %s must begin with %s followed by a bucket name", url, s3Prefix)
	}
	if !strings.Contains(url[len(s3Prefix):], "/") {
		return url[len(s3Prefix):], "", nil
	}
	firstSlash := strings.Index(url[len(s3Prefix):], "/") + len(s3Prefix)
	return url[len(s3Prefix):firstSlash], url[firstSlash + 1:], nil
}
