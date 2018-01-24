package aws

import (
	"testing"
)

func testBucketKeyFromURL (t *testing.T) {
	url := "s3://my-bucket/my-dir/my-image.png"
	bucket, key, err := bucketKeyFromURL(url)
	if bucket != "my-bucket" || key != "my-dir/my-image.png" || err != nil {
		t.Fatal(url, bucket, key, err)
	}
	url = "s3://"
	bucket, key, err = bucketKeyFromURL(url)
	if err == nil {
		t.Fatal(url, err)
	}
}
