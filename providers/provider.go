package providers

import (
	"github.com/nathanwilk7/zcloud/storage"
	"github.com/nathanwilk7/zcloud/compute"
)

type Provider interface {
	storage.StorageProvider
	compute.ComputeProvider
}

// TODO: funcs that give a provider back
