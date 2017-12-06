package main

type StorageProvider interface {
	Cp (src, dst string)
}
