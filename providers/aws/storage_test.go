package aws

import (
	"testing"
)

func testGetS3BucketKeyFromURL (t *testing.T) {
	url := "s3://my-bucket/my-dir/my-image.png"
	bucket, key, err := getS3BucketKeyFromURL(url)
	if bucket != "my-bucket" || key != "my-dir/my-image.png" || err != nil {
		t.Fatal(url, bucket, key, err)
	}
	url = "s3://"
	bucket, key, err = getS3BucketKeyFromURL(url)
	if err == nil {
		t.Fatal(url, err)
	}
}
