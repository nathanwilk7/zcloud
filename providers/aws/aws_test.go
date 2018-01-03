package aws

import (
	"testing"
)

func TestConvertURL (t *testing.T) {
	url := "cloud://bucket/file.txt"
	if res := convertURL(url); res != "s3://bucket/file.txt" {
		t.Fatal(res)
	}
}
