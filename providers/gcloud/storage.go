package gcloud

import (
	"fmt"
	"io"
	"os"
	"strings"
	
	"github.com/nathanwilk7/zcloud/storage"
	"google.golang.org/api/iterator"

	gs "cloud.google.com/go/storage"
)

func (p gcloudProvider) Upload (params storage.UploadParams) (string, error) {
	f, err := os.Open(params.Src)
	if err != nil {
		return "", err
	}
	defer f.Close()

	bucket, object, err := bucketObjectFromURL(params.Dest)
	if err != nil {
		return "", err
	}
	ctx := getContext()
	client, err := getClient(ctx)
	if err != nil {
		return "", err
	}
	wc := client.Bucket(bucket).Object(object).NewWriter(ctx)
	if _, err = io.Copy(wc, f); err != nil {
		return "", err
	}
	if err := wc.Close(); err != nil {
		return "", err
	}
	return fmt.Sprintf("Uploaded file %s", params.Src), nil
}

func (p gcloudProvider) Download (params storage.DownloadParams) (string, error) {
	bucket, object, err := bucketObjectFromURL(params.Src)
	if err != nil {
		return "", err
	}
	ctx := getContext()
	client, err := getClient(ctx)
	if err != nil {
		return "", err
	}
	r, err := client.Bucket(bucket).Object(object).NewReader(ctx)
	if err != nil {
		return "", err
	}
	defer r.Close()
	
	f, err := os.Create(params.Dest)
	if err != nil {
		return "", err
	}
	defer f.Close()
	
	if _, err = io.Copy(f, r); err != nil {
		return "", err
	}
	if err := r.Close(); err != nil {
		return "", err
	}
	return fmt.Sprintf("Downloaded file %s", params.Src), nil
}

func (p gcloudProvider) Ls (params storage.LsParams) (string, error) {
	ctx := getContext()
	client, err := getClient(ctx)
	if err != nil {
		return "", err
	}
	bucket, object, err := bucketObjectFromURL(params.Url)
	if err != nil {
		return "", err
	}
	query := &gs.Query{
		Prefix: object,
	}
	objects := client.Bucket(bucket).Objects(ctx, query)
	out := []string{}
	for {
		object, err := objects.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return "", err
		}
		out = append(out, object.Name)
	}
	return strings.Join(out, "\n"), nil
}

func (p gcloudProvider) Rm (params storage.RmParams) (string, error) {
	return "", nil
}

func (p gcloudProvider) Mv (params storage.MvParams) (string, error) {
	return "", nil
}

func (p gcloudProvider) Mb (params storage.MbParams) (string, error) {
	return "", nil
}

func (p gcloudProvider) Rb (params storage.RmParams) (string, error) {
	return "", nil
}

func (p gcloudProvider) Sync (params storage.SyncParams) (string, error) {
	return "", nil
}

const gsStr = "gs"
const gsPrefix = gsStr + "://"

func (p gcloudProvider) StorageURLPrefixReplacement() string {
	return gsStr
}

func bucketObjectFromURL (url string) (string, string, error) {
	if len(url) <= len(gsPrefix) {
		return "", "", fmt.Errorf("Converted URL %s must begin with %s followed by a bucket name", url, gsPrefix)
	}
	if !strings.Contains(url[len(gsPrefix):], "/") {
		return url[len(gsPrefix):], "", nil
	}
	firstSlash := strings.Index(url[len(gsPrefix):], "/") + len(gsPrefix)
	return url[len(gsPrefix):firstSlash], url[firstSlash + 1:], nil
}

