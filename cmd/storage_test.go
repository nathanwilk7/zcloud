package cmd

import (
	"testing"

	"os"
)

func TestMustGetStorageProvider (t *testing.T) {
	_ = mustGetStorageProvider(os.Getenv(storageProvEnv), os.Getenv(provEnv))
}
