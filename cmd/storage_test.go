package cmd

import (
	"testing"
)

func TestMustGetStorageProvider (t *testing.T) {
	_ = mustGetStorageProvider()
}
