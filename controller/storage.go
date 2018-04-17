package controller

import (
	"path/filepath"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
	
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
		oqp := &z.ObjectsQueryParams{
			Prefix: key,
		}
		objects, err := b.ObjectsQuery(oqp)
		if err != nil {
			return 0, err
		}
		keys = keysFromObjects(objects)
		fps = fpsFromFilepathKeys(fp, keys)
	}
	for i := range fps {
		err := downloadFile(b, keys[i], fps[i], false)
		if err != nil {
			return i, err
		}
	}
	return len(fps), nil
}

func downloadFile (b z.Bucket, k, fpath string, deltaTransfer bool) error {
	d := filepath.Dir(fpath)
	err := os.MkdirAll(d, os.ModePerm)
	if err != nil {
		return err
	}
	f, err := os.Create(fpath)
	if err != nil {
		return err
	}
	fi, err := f.Stat()
	if err != nil {
		return err
	}
	o := b.Object(k)
	oi, err := o.Info()
	if err != nil {
		return err
	}
	if deltaTransfer &&
		!shouldReplace(oi.Size(), int(fi.Size()), oi.LastModified(), fi.ModTime()) {
		return nil
	}
	r, err := o.Reader()
	if err != nil {
		f.Close()
		return err
	}
	_, err = io.Copy(f, r)
	if err != nil {
		r.Close()
		f.Close()
		return err
	}
	err = r.Close()
	if err != nil {
		f.Close()
		return err
	}
	err = f.Close()
	if err != nil {
		return err
	}
	return nil
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
		err := uploadFile(fps[i], b, keys[i], false)
		if err != nil {
			return i, err
		}
	}
	return len(fps), nil
}

func uploadFile (path string, b z.Bucket, k string, deltaTransfer bool) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	fi, err := f.Stat()
	if err != nil {
		return err
	}
	o := b.Object(k)
	oi, err := o.Info()
	// don't return err if the object doesn't exist
	switch err.(type) {
	case z.ErrObjectDoesNotExist:
		// placeholder text
	default:
		return err
	}
	switch err.(type) {
	case z.ErrObjectDoesNotExist:
		// placeholder text
	default:
		if deltaTransfer &&
			!shouldReplace(int(fi.Size()), oi.Size(), fi.ModTime(), oi.LastModified()) {
			return nil
		}
	}
	w, err := o.Writer()
	if err != nil {
		f.Close()
		return err
	}
	_, err = io.Copy(w, f)
	if err != nil {
		w.Close()
		f.Close()
		return err
	}
	err = w.Close()
	if err != nil {
		f.Close()
		return err
	}
	err = f.Close()
	if err != nil {
		return err
	}
	return nil
}

func keysFromFilepaths (fps []string, fileprefix, urlprefix string) []string {
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
	// TODO: Why is this here?
	if urlprefix[len(urlprefix) - 1] == '/' {
		urlprefix += "/"
	}
	urlprefix = ensureTrailingSlash(urlprefix)
	return fmt.Sprintf("%s%s", urlprefix, fp)
}

func recursiveFilepaths (fp string) []string {
	fps := []string{}
	filepath.Walk(fp, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
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
		prefix = ensureTrailingSlash(prefix)
	}
	postfix := strings.Replace(key, prefix, "", 1)
	i := strings.Index(postfix, "/")
	if i != -1 {
		return postfix[:i + 1]
	}
	return postfix
}

func ensureTrailingSlash (s string) string {
	if len(s) == 0 {
		return "/"
	}
	if s[len(s) - 1] != '/' {
		s += "/"
	}
	return s
}

func ensureTrailingSlashIfDir (filepath string) (string, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return filepath, err
	}
	fi, err := f.Stat()
	if err != nil {
		return filepath, err
	}
	if fi.IsDir() {
		filepath = ensureTrailingSlash(filepath)
	}
	return filepath, nil
}

type RmParams struct {
	Url string
	Recursive bool
}

func Rm (pp ProvParams, rp RmParams, o out.Out) {
	p, err := z.NewProvider(zppFromPp(pp))
	if err != nil {
		o.Fatal(err)
	}
	if !isCloudURL(rp.Url) {
		o.Fatalf("%v is not a valid zCloud URL", rp.Url)
	}
	bn, k, err := bucketNameKey(rp.Url)
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
	if rp.Recursive {
		for i := range objects {
			err = objects[i].Delete()
			if err != nil {
				o.Fatal(err)
			}
		}
		o.Messagef("%v objects were deleted from %v\n", len(objects), rp.Url)
	} else {
		if len(objects) != 1 {
			o.Fatalf("Object doesn't exist or more than one object was returned: %v", rp.Url)
		}
		objects[0].Delete()
		o.Messagef("%v was deleted\n", rp.Url)
	}
}

type SyncParams struct {
	Src, Dest string
}

func Sync (pp ProvParams, sp SyncParams, o out.Out) {
	p, err := z.NewProvider(zppFromPp(pp))
	if err != nil {
		o.Fatal(err)
	}
	var sbn, sk, dbn, dk string
	if isCloudURL(sp.Src) {
		sbn, sk, err = bucketNameKey(sp.Src)
		if err != nil {
			o.Fatal(err)
		}
	}
	if isCloudURL(sp.Dest) {
		dbn, dk, err = bucketNameKey(sp.Dest)
		if err != nil {
			o.Fatal(err)
		}
	}
	if isCloudURL(sp.Src) && isCloudURL(sp.Dest) {
		err = syncCloudToCloud(p, sbn, sk, dbn, dk)
	} else if !isCloudURL(sp.Src) && isCloudURL(sp.Dest) {
		err = syncLocalToCloud(p, sp.Src, dbn, dk)
	} else if isCloudURL(sp.Src) && !isCloudURL(sp.Dest) {
		err = syncCloudToLocal(p, sbn, sk, sp.Dest)
	} else {
		o.Fatalf("Src or dest must be a cloud url, src: %v, dest: %v", sp.Src, sp.Dest)
	}
	if err != nil {
		o.Fatal(err)
	}
}

func syncCloudToCloud (p z.Provider, sbn, sk, dbn, dk string) error {
	soq := &z.ObjectsQueryParams{
		Prefix: sk,
	}
	sos, err := p.Bucket(sbn).ObjectsQuery(soq)
	if err != nil {
		return err
	}
	db := p.Bucket(dbn)
	for _, so := range sos {
		dk := dk + strings.Replace(so.Key(), sk, "", 1)
		do := db.Object(dk)
		soi, err := so.Info()
		if err != nil {
			return err
		}
		doi, err := do.Info()
		if err == nil &&
			!shouldReplace(soi.Size(), doi.Size(), soi.LastModified(), doi.LastModified()) {
			continue
		}
		err = so.CopyTo(do)
		if err != nil {
			return err
		}
	}
	return nil
}

func syncLocalToCloud (p z.Provider, fileprefix, dbn, dk string) error {
	db := p.Bucket(dbn)
	fileprefix, err := ensureTrailingSlashIfDir(fileprefix)
	if err != nil {
		return err
	}
	return filepath.Walk(fileprefix, func (path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			newKey := keyFromFilepath(path, fileprefix, dk)
			err = uploadFile(path, db, newKey, true)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func syncCloudToLocal (p z.Provider, bn, k, fileprefix string) error {
	oq := &z.ObjectsQueryParams{
		Prefix: k,
	}
	b := p.Bucket(bn)
	obs, err := b.ObjectsQuery(oq)
	if err != nil {
		return err
	}
	fileprefix, err = ensureTrailingSlashIfDir(fileprefix)
	if err != nil {
		return err
	}
	for _, o := range obs {
		filepath := fileprefix + strings.Replace(o.Key(), k, "", 1)
		err = downloadFile(b, o.Key(), filepath, true)
		if err != nil {
			return err
		}
	}
	return nil
}

type TransferParams struct {
	Src, Dest string
	DestProv string
	Recursive bool
}

func Transfer (pp ProvParams, tp TransferParams, o out.Out) {
	srcProv := pp.Name
	sp, err := z.NewProvider(zppFromPp(pp))
	pp.Name = tp.DestProv
	dp, err := z.NewProvider(zppFromPp(pp))
	sbn, sk, err := bucketNameKey(tp.Src)
	if err != nil {
		o.Fatal(err)
	}
	dbn, dk, err := bucketNameKey(tp.Dest)
	if err != nil {
		o.Fatal(err)
	}
	sb := sp.Bucket(sbn)
	var sos []z.Object
	if tp.Recursive || sk == "" {
		oqp := &z.ObjectsQueryParams{
			Prefix: sk,
		}
		sos, err = sb.ObjectsQuery(oqp)
		if err != nil {
			o.Fatal(err)
		}
	} else {
		so := sb.Object(sk)
		sos = []z.Object{so}
	}
	db := dp.Bucket(dbn)
	for _, so := range sos {
		do := db.Object(dk + strings.Replace(so.Key(), sk, "", 1))
		r, err := so.Reader()
		if err != nil {
			o.Fatal(err)
		}
		w, err := do.Writer()
		if err != nil {
			o.Fatal(err)
		}
		_, err = io.Copy(w, r)
		if err != nil {
			o.Fatal(err)
		}
		err = r.Close()
		if err != nil {
			o.Fatal(err)
		}
		err = w.Close()
		if err != nil {
			o.Fatal(err)
		}
	}
	o.Messagef("Transferred %v files from %v to %v\n", len(sos), srcProv, tp.DestProv)
}

func shouldReplace (ss, ds int, slm, dlm time.Time) bool {
	return ss != ds ||
		dlm.Before(slm)
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
