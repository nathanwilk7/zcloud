package controller

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBaseNameOrDir (t *testing.T) {
	k := "a.txt"
	p := ""
	s := firstPathEl(k, p)
	if s != k {
		t.Fatalf("f: %v k: %v s: %v", k, p, s)
	}
	k = "dir/subdir/a.txt"
	p = ""
	s = firstPathEl(k, p)
	if s != "dir/" {
		t.Fatalf("f: %v k: %v s: %v", k, p, s)
	}
	k = "dir/subdir/a.txt"
	p = "dir/subdir"
	s = firstPathEl(k, p)
	if s != "a.txt" {
		t.Fatalf("f: %v k: %v s: %v", k, p, s)
	}
	k = "dir/subdir/a.txt"
	p = "dir/subdir/"
	s = firstPathEl(k, p)
	if s != "a.txt" {
		t.Fatalf("f: %v k: %v s: %v", k, p, s)
	}
	k = "dir/subdir/a.txt"
	p = "dir"
	s = firstPathEl(k, p)
	if s != "subdir/" {
		t.Fatalf("f: %v k: %v s: %v", k, p, s)
	}
	k = "dir/subdir/a.txt"
	p = "dir/"
	s = firstPathEl(k, p)
	if s != "subdir/" {
		t.Fatalf("f: %v k: %v s: %v", k, p, s)
	}
}

func TestRecursiveFilepaths (t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	var (
		testDirA = wd + "/testdata"
		testDirB = testDirA + "/dir"
		testDirC = testDirB + "/subdir"
		testFileA = testDirA + "/a.txt"
		testFileB = testDirB + "/b.txt"
		testFileC = testDirC + "/c.txt"
	)
	err = os.Mkdir(testDirA, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
	_, err = os.Create(testFileA)
	if err != nil {
		t.Fatal(err)
	}
	err = os.Mkdir(testDirB, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
	_, err = os.Create(testFileB)
	if err != nil {
		t.Fatal(err)
	}
	err = os.Mkdir(testDirC, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
	_, err = os.Create(testFileC)
	if err != nil {
		t.Fatal(err)
	}
	fps := recursiveFilepaths(testDirA)
	assert := assert.New(t)
	assert.Contains(fps, testFileA)
	assert.Contains(fps, testFileB)
	assert.Contains(fps, testFileC)
	fps = recursiveFilepaths(testDirB)
	assert.Contains(fps, testFileB)
	assert.Contains(fps, testFileC)
	fps = recursiveFilepaths(testDirC)
	assert.Contains(fps, testFileC)
	os.RemoveAll(testDirA)
}

func TestKeyFromFilepath (t *testing.T) {
	fp := "/base/a.txt"
	fileprefix := "/"
	urlprefix := "dir"
	r := keyFromFilepath(fp, fileprefix, urlprefix)
	if r != "dir/base/a.txt" {
		t.Fatalf("fp %v fileprefix %v urlprefix %v r %v", fp, fileprefix, urlprefix, r)
	}
	fp = "/base/a.txt"
	fileprefix = "/base"
	urlprefix = "dir"
	r = keyFromFilepath(fp, fileprefix, urlprefix)
	if r != "dir/a.txt" {
		t.Fatalf("fp %v fileprefix %v urlprefix %v r %v", fp, fileprefix, urlprefix, r)
	}
	fp = "/base/a.txt"
	fileprefix = "/base"
	urlprefix = ""
	r = keyFromFilepath(fp, fileprefix, urlprefix)
	if r != "a.txt" {
		t.Fatalf("fp %v fileprefix %v urlprefix %v r %v", fp, fileprefix, urlprefix, r)
	}
}

func TestFpsFromFilepathKeys (t *testing.T) {
	keys := []string{"a.txt", "dir/b.txt", "dir/subdir/c.txt"}
	fp := "/home"
	fps := fpsFromFilepathKeys(fp, keys)
	if fps[0] != "/home/a.txt" {
		t.Fatal("keys %v fp %v fps %v", keys, fp, fps)
	}
	if fps[1] != "/home/dir/b.txt" {
		t.Fatal("keys %v fp %v fps %v", keys, fp, fps)
	}
	if fps[2] != "/home/dir/subdir/c.txt" {
		t.Fatal("keys %v fp %v fps %v", keys, fp, fps)
	}
}
