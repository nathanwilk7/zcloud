package controller

import (
	"path/filepath"
	"fmt"
	"io"
	"os"
	"strings"
	
	"github.com/nathanwilk7/zcloud/out"
	z "github.com/nathanwilk7/zcloud-go"

	"github.com/nathanwilk7/stringset"
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
	url, fp, err := getURLFilepath(cp.Src, cp.Dest)
	if err != nil {
		o.Fatal(err)
	}
	fp, err = filepath.Abs(fp)
	if err != nil {
		o.Fatal(err)
	}
	bn, k, err := bucketNameKey(url)
	if err != nil {
		o.Fatal(err)
	}
	b := p.Bucket(bn)
	var n int
	if isCloudURL(cp.Src) {
		n, err = Download(p, b, fp, k, cp.Recursive)
	} else {
		n, err = Upload(p, b, fp, k, cp.Recursive)
	}
	if err != nil {
		o.Fatal(err)
	}
	if !cp.Recursive {
		o.Messagef("Successfully copied %v to %v\n", cp.Src, cp.Dest)
	} else {
		o.Messagef("Successfully copied %v files to %v\n", n, cp.Dest)
	}
}

func Download (p z.Provider, b z.Bucket, fp, key string, recursive bool) (int, error) {
	var fps []string
	var keys []string
	if !recursive {
		fps = []string{fp}
		keys = []string{key}
	} else {
		oqp := z.ObjectsQueryParams{
			Prefix: key,
		}
		objects, err := b.ObjectsQuery(&oqp)
		if err != nil {
			return 0, err
		}
		keys = keysFromObjects(objects)
		fps = fpsFromFilepathKeys(fp, keys)
	}
	for i := range fps {
		d := filepath.Dir(fps[i])
		err := os.MkdirAll(d, os.ModePerm)
		if err != nil {
			return 0, err
		}
		f, err := os.Create(fps[i])
		if err != nil {
			return 0, err
		}
		r, err := b.Object(keys[i]).Reader()
		if err != nil {
			f.Close()
			return 0, err
		}
		_, err = io.Copy(f, r)
		if err != nil {
			r.Close()
			f.Close()
			return 0, err
		}
		err = r.Close()
		if err != nil {
			f.Close()
			return 0, err
		}
		err = f.Close()
		if err != nil {
			return 0, err
		}
	}
	return len(fps), nil
}

func keysFromObjects (objects []z.Object) []string {
	keys := make([]string, len(objects))
	for i := range objects {
		keys[i] = objects[i].Key()
	}
	return keys
}

func fpsFromFilepathKeys(fp string, keys []string) []string {
	fps := make([]string, len(keys))
	for i := range keys {
		fps[i] = filepath.Join(fp, keys[i])
	}
	return fps
}

func Upload (p z.Provider, b z.Bucket, fp, key string, recursive bool) (int, error) {
	var fps []string
	var keys []string
	if !recursive {
		fps = []string{fp}
		keys = []string{key}
	} else {
		fps = recursiveFilepaths(fp)
		keys = keysFromFilepaths(fps, fp, key)
	}
	for i := range fps {
		f, err := os.Open(fps[i])
		if err != nil {
			return 0, err
		}
		w, err := b.Object(keys[i]).Writer()
		if err != nil {
			f.Close()
			return 0, err
		}
		_, err = io.Copy(w, f)
		if err != nil {
			w.Close()
			f.Close()
			return 0, err
		}
		err = w.Close()
		if err != nil {
			f.Close()
			return 0, err
		}
		err = f.Close()
		if err != nil {
			return 0, err
		}
	}
	return len(fps), nil
}

func keysFromFilepaths(fps []string, fileprefix, urlprefix string) []string {
	keys := make([]string, len(fps))
	for i := range fps {
		keys[i] = keyFromFilepath(fps[i], fileprefix, urlprefix)
	}
	return keys
}

func keyFromFilepath (fp, fileprefix, urlprefix string) string {
	fp = strings.Replace(fp, fileprefix, "", 1)
	if len(fp) == 0 {
		return ""
	}
	if fp[0] == '/' {
		fp = fp[1:]
	}
	if len(urlprefix) == 0 {
		return fp
	}
	if urlprefix[len(urlprefix) - 1] == '/' {
		urlprefix += "/"
	}
	urlprefix = maybeAppendSlash(urlprefix)
	return fmt.Sprintf("%s%s", urlprefix, fp)
}

func recursiveFilepaths (fp string) []string {
	fps := []string{}
	filepath.Walk(fp, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			fps = append(fps, path)
		}
		return nil
	})
	return fps
}

func getURLFilepath (a string, b string) (string, string, error) {
	var url, filepath string
	if isCloudURL(a) && !isCloudURL(b) {
		url = a
		filepath = b
	} else if !isCloudURL(a) && isCloudURL(b) {
		filepath = a
		url = b
	} else {
		err := fmt.Errorf(
			"Exactly one of the source and cpdDestination url's must be a cloud url with the format cloud://...: %s, %s",
			a,
			b,
		)
		return "", "", err
	}
	return url, filepath, nil
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
	objects, err := p.Bucket(bn).ObjectsQuery(q)
	if err != nil {
		o.Fatal(err)
	}
	ss := stringset.New()
	if ls.Recursive {
		for i := range objects {
			ss.Add(objects[i].Key())
		}
	} else {
		for i := range objects {
			ss.Add(firstPathEl(objects[i].Key(), k))
		}
	}
	ssSlice := ss.ToSlice()
	fis := make([]out.FileInfo, len(ssSlice))
	for i, s := range ssSlice {
		fis[i].Name = s
	}
	o.ListFileInfos(fis)
}

func firstPathEl (key, prefix string) string {
	if len(prefix) != 0 {
		prefix = maybeAppendSlash(prefix)
	}
	postfix := strings.Replace(key, prefix, "", 1)
	i := strings.Index(postfix, "/")
	if i != -1 {
		return postfix[:i + 1]
	}
	return postfix
}

func maybeAppendSlash (s string) string {
	if len(s) == 0 {
		return "/"
	}
	if s[len(s) - 1] != '/' {
		s += "/"
	}
	return s
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
