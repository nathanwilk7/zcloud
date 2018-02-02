package aws

import (
	"os"
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

func TestGetFbksToUpload (t *testing.T) {
	os.Mkdir("testdata", os.ModePerm)
	defer os.Remove("testdata")
	f, err := os.Create("test-fbks.txt")
	defer f.Close()
	defer os.Remove("testdata/test-fbks.txt")
	if err != nil {
		t.Fatal("couldn't create file")
	}
	_, err = f.Write([]byte("fbk"))
	if err != nil {
		t.Fatal("couldn't write to file")
	}
	src := "testdata"
	bucket := "bucket"
	key := "key"
	recursive := false
	fbks, err := getFbksToUpload(src, bucket, key, recursive)
	if err != nil {
		t.Fatal(err)
	}
	if len(fbks) != 1 && fbks[0].File.Name() != "test-fbks.txt" {
		t.Fatal("Didn't get the right fbk's back")
	}
}

func TestCreateAndAppendFbk (t *testing.T) {
	fbks := []FileBucketKey{}
	os.Mkdir("testdata", os.ModePerm)
	defer os.Remove("testdata")
	filename, bucket, key := "testdata/test-create.txt", "bucket", "key"
	var err error
	fbks, err = createAndAppendFbk(fbks, filename, bucket, key)
	f, err := os.Open("testdata/test-create.txt")
	defer f.Close()
	defer os.Remove("testdata/test-create.txt")
	if err != nil {
		t.Fatal("could not open file")
	}
}
