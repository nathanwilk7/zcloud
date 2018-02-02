package aws

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/nathanwilk7/zcloud/storage"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/nathanwilk7/stringset"
)

func (p awsProvider) Upload (params storage.UploadParams) (string, error) {
	bucket, key, err := bucketKeyFromURL(params.Dest)
	if err != nil {
		return "", err
	}
	// hold onto the fbks because we'll need them to close the files after we upload
	fbks, err := getFbksToUpload(params.Src, bucket, key, params.Recursive)
	if err != nil {
		return "", err
	}
	if len(fbks) == 0 {
		return "", fmt.Errorf("No objects to upload were specified by source: %s", params.Src)
	}
	err = uploadObjects(p.Session, fbks)
	if err != nil {
		return "", err
	}
	for _, fbk := range fbks {
		fbk.File.Close()
	}
	return "files uploaded", nil
}

type FileBucketKey struct {
	File *os.File
	Bucket, Key string
}

// TODO: test on empty dir
func getFbksToUpload (src, bucket, key string, recursive bool) ([]FileBucketKey, error) {
	fbks := []FileBucketKey{}
	if !recursive {
		var err error
		fbks, err = appendToFbks(fbks, src, bucket, key)
		if err != nil {
			return fbks, err
		}
	} else {
		err := filepath.Walk(src, func (path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}
			// get relative filepath to use as second half of key. This cuts off src from path because we're guaranteed that path contains
			// src.
			relpath, err := filepath.Rel(src, path)
			if err != nil {
				return err
			}
			// if relative filepath is just current dir, use the filename instead
			if relpath == "." {
				relpath = filepath.Base(src)
			}
			// if the key exists, prepend it to the relative path. This allows files to be uploaded into a directory on aws instead of
			// only at the bucket level
			var formattedKey string
			if key == "" {
				formattedKey = relpath
			} else {
				formattedKey = key + "/" + relpath
			}
			fbks, err = appendToFbks(fbks, path, bucket, formattedKey)
			if err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			return fbks, fmt.Errorf("error occured while walking directories: %v", err)
		}
	}
	return fbks, nil
}

func appendToFbks (fbks []FileBucketKey, filepath, bucket, key string) ([]FileBucketKey, error) {
	f, err  := os.Open(filepath)
	if err != nil {
		return fbks, fmt.Errorf("failed to open file %s, %v", filepath, err)
	}
	fbks = append(fbks, FileBucketKey{
		File: f,
		Bucket: bucket,
		Key: key,
	})
	return fbks, nil
}

func uploadObjects(sess *session.Session, fbks []FileBucketKey) error {
	objects := []s3manager.BatchUploadObject{}
	for _, fbk := range fbks {
		objects = append(objects, s3manager.BatchUploadObject{
			Object: &s3manager.UploadInput{
				Bucket: aws.String(fbk.Bucket),
				Key:    aws.String(fbk.Key),
				Body:   fbk.File,
			},
		})
	}
	iter := &s3manager.UploadObjectsIterator{Objects: objects}
	uploader := s3manager.NewUploader(sess)
	if err := uploader.UploadWithIterator(aws.BackgroundContext(), iter); err != nil {
		return err
	}
	return nil
}

func (p awsProvider) Download (params storage.DownloadParams) (string, error) {
	svc := s3.New(p.Session)
	downloader := s3manager.NewDownloader(p.Session)
	bucket, key, err := bucketKeyFromURL(params.Src)
	if err != nil {
		return "", err
	}
	if !params.Recursive {
		f, err := os.Create(params.Dest)
		if err != nil {
			return "", fmt.Errorf("failed to create file %q, %v", params.Dest, err)
		}
		defer f.Close()
		
		_, err = downloader.Download(f, &s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		})
		if err != nil {
			return "", fmt.Errorf("failed to download file %v", err)
		}
		return fmt.Sprintf("download: %s to %s", params.Src, params.Dest), nil
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
	fbks := []FileBucketKey{}
	for _, content := range result.Contents {
		// get local filename from dest and object key
		filename := params.Dest + *content.Key
		// if we're downloading from a cloud subdirectory, then remove it from the filepath so we don't create unnecessary subdirs
		if key != "" {
			filename += strings.Replace(filename, key, "", 1)
		}
		os.MkdirAll(filepath.Dir(filename), os.ModePerm)
		// Don't create a file for directories (which have a Size of 0)
		if *content.Size == 0 {
			continue
		}
		fbks, err = createAndAppendFbk (fbks, filename, bucket, *content.Key)
		if err != nil {
			return "", err
		}
	}
	err = downloadFbks(fbks, downloader)
	if err != nil {
		return "", err
	}
	for _, fbk := range fbks {
		fbk.File.Close()
	}
	return fmt.Sprintf("downloaded files"), nil
}

func createAndAppendFbk (fbks []FileBucketKey, filename, bucket, key string) ([]FileBucketKey, error) {
	f, err := os.Create(filename)
	if err != nil {
		return fbks, fmt.Errorf("failed to create file %q, %v", filename, err)
	}
	fbks = append(fbks, FileBucketKey {
		Bucket: bucket,
		Key: key,
		File: f,
	})
	return fbks, nil
}

func downloadFbks (fbks []FileBucketKey, downloader* s3manager.Downloader) error {
	objects := []s3manager.BatchDownloadObject{}
	for _, fbk := range fbks {
		objects = append(objects, s3manager.BatchDownloadObject {
			Object: &s3.GetObjectInput {
				Bucket: aws.String(fbk.Bucket),
				Key: aws.String(fbk.Key),
			},
			Writer: fbk.File,
		})
	}
	iter := &s3manager.DownloadObjectsIterator{Objects: objects}
	if err := downloader.DownloadWithIterator(aws.BackgroundContext(), iter); err != nil {
		return err
	}
	return nil
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
	ss := stringset.NewStringset()
	for _, content := range result.Contents {
		key := *content.Key
		if strings.Contains(key, "/") && !params.Recursive {
			ss.Add(key[:strings.Index(key, "/") + 1])
		} else {
			ss.Add(key)
		}		
	}
	res := strings.Join(ss.ToSlice(), "\n")
	return res, nil
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
