package controller

import (
	"fmt"
	"io"
	"os"
	"strings"
	
	"github.com/nathanwilk7/zcloud/out"
	z "github.com/nathanwilk7/zcloud-go"

	"github.com/nathanwilk7/zcloud/stringset"
)
		

type ProvParams struct {
	Name string
	AwsId, AwsSecret, AwsRegion string
	GCloudProjectID string
}

type CpParams struct {
	Src, Dest string
	Recursive bool
}

func Cp (pp ProvParams, cp CpParams, o out.Out) {
	p, err := z.NewProvider(zppFromPp(pp))
	if err != nil {
		o.Fatal(err)
	}
	url, filename, err := getURLFilename(cp.Src, cp.Dest)
	if err != nil {
		o.Fatal(err)
	}
	bn, k, err := bucketNameKey(url)
	if err != nil {
		o.Fatal(err)
	}
	f, err := os.Create(filename)
	if err != nil {
		o.Fatal(err)
	}
	r, err := p.Bucket(bn).Object(k).Reader()
	if err != nil {
		r.Close()
		o.Fatal(err)
	}
	_, err = io.Copy(f, r)
	if err != nil {
		r.Close()
		o.Fatal(err)
	}
	err = r.Close()
	if err != nil {
		o.Fatal(err)
	}
	o.Messagef("Successfully copied %v to %v\n", cp.Src, cp.Dest)
}

func getURLFilename (a string, b string) (string, string, error) {
	var url, filename string
	if isCloudURL(a) && !isCloudURL(b) {
		url = a
		filename = b
	} else if !isCloudURL(a) && isCloudURL(b) {
		filename = a
		url = b
	} else {
		err := fmt.Errorf(
			"Exactly one of the source and cpdDestination url's must be a cloud url with the format cloud://...: %s, %s",
			a,
			b,
		)
		return "", "", err
	}
	return url, filename, nil
}

type LsParams struct {
	Url string
	Recursive bool
}

func Ls (pp ProvParams, ls LsParams, o out.Out) {
	p, err := z.NewProvider(zppFromPp(pp))
	if err != nil {
		o.Fatal(err)
	}
	if !isCloudURL(ls.Url) {
		o.Fatalf("%v is not a valid zCloud URL", ls.Url)
	}
	bn, k, err := bucketNameKey(ls.Url)
	if err != nil {
		o.Fatal(err)
	}
	q := &z.ObjectsQueryParams{
		Prefix: k,
	}
	os, err := p.Bucket(bn).ObjectsQuery(q)
	if err != nil {
		o.Fatal(err)
	}
	ss := stringset.New()
	if ls.Recursive {
		for i := range os {
			ss.Add(os.Name())
		}
	} else {
		for i := range os {
			
		}
	}
	fis[i].Name = os[i].Key()
	fis := make([]out.FileInfo, len(os))
	o.ListFileInfos(fis)
}

func baseNameOrDir (filename, key string) {
	
}

func zppFromPp (pp ProvParams) z.ProviderParams {
	return z.ProviderParams{
		Name: pp.Name,
		AwsId: pp.AwsId,
		AwsSecret: pp.AwsSecret,
		AwsRegion: pp.AwsRegion,
		GCloudProjectID: pp.GCloudProjectID,
	}
}

const (
	cloudStr = "cloud"
	cloudURLPrefix = cloudStr + "://"
)

func convertURL (url, replacement string) string {
	if isCloudURL(url) {
		return strings.Replace(url, cloudStr, replacement, 1)
	}
	return url
}

func isCloudURL (url string) bool {
	if len(url) > len(cloudURLPrefix) {
		if url[:len(cloudURLPrefix)] == cloudURLPrefix {
			return true
		}
	}
	return false
}

func bucketNameKey (url string) (string, string, error) {
	if len(url) <= len(cloudURLPrefix) {
		return "", "", fmt.Errorf("Converted URL %s must begin with %s followed by a bucket name", url, cloudURLPrefix)
	}
	if !strings.Contains(url[len(cloudURLPrefix):], "/") {
		return url[len(cloudURLPrefix):], "", nil
	}
	firstSlash := strings.Index(url[len(cloudURLPrefix):], "/") + len(cloudURLPrefix)
	return url[len(cloudURLPrefix):firstSlash], url[firstSlash + 1:], nil
}
