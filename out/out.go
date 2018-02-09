package out

import (
	"fmt"
	"log"
)

func New () Out {
	return DefaultOut{}
}

type Out interface {
	Fatal (is ...interface{})
	Fatalf (format string, v ...interface{})
	Messageln (msg string)
	Messagef (format string, v ...interface{})
	ListFileInfos (fis []FileInfo)
}

type DefaultOut struct {}

type FileInfo struct {
	Name string
}

func (o DefaultOut) Fatal (v ...interface{}) {
	log.Fatal(v)
}

func (o DefaultOut) Fatalf (format string, v ...interface{}) {
	log.Fatalf(format, v)
}

func (o DefaultOut) Messageln (msg string) {
	fmt.Println(msg)
}

func (o DefaultOut) Messagef (format string, v ...interface{}) {
	fmt.Printf(format, v)
}

func (o DefaultOut) ListFileInfos (fis []FileInfo) {
	fmt.Println(fis)
}
