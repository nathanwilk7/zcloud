package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/nathanwilk7/zcloud/storage"

	"github.com/spf13/cobra"
)

var StorageCmd = &cobra.Command{
	Use:   "storage",
	Short: "Blob storage",
	Long:  "Perform various operations related to blob storage",
}

var cpRecursive bool
var lsRecursive bool

func init () {
	CpCmd.Flags().BoolVarP(&cpRecursive, "recursive", "r", false, "Recursively copy from src")
	StorageCmd.AddCommand(CpCmd)

	LsCmd.Flags().BoolVarP(&lsRecursive, "recursive", "r", false, "Recursively list")
	StorageCmd.AddCommand(LsCmd)

	RootCmd.AddCommand(StorageCmd)
}

const storageProvEnv = "ZCLOUD_STORAGE_PROV"

func getStorageProvider (storageProvider string, provider string) (storage.StorageProvider, error) {
	if p, ok := storageProviders[storageProvider]; ok {
		return p, nil
	}
	if p, err := getProvider(storageProvider); err == nil {
		return p, nil
	}
	if p, err := getProvider(provider); err == nil {
		return p, nil
	}
	return nil, fmt.Errorf("%s and %s were not valid StorageProviders", storageProvider, provider)
}

var storageProviders map[string]storage.StorageProvider = map[string]storage.StorageProvider {}


var CpCmd = &cobra.Command{
	Use:   "cp",
	Short: "Copy objects",
	Long:  "Copy objects to/from a provider",
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		src, dest := args[0], args[1]
		var msg string
		var err error
		p := mustGetStorageProvider(os.Getenv(storageProvEnv), os.Getenv(provEnv))
		replacement := p.StorageURLPrefixReplacement()
		if isCloudURL(src) && !isCloudURL(dest) {
			params := storage.NewDownloadParams(convertURL(src, replacement), dest)
			params.Recursive = cpRecursive
			msg, err = p.Download(params)
		} else if !isCloudURL(src) && isCloudURL(dest) {
			params := storage.NewUploadParams(src, convertURL(dest, replacement))
			params.Recursive = cpRecursive
			msg, err = p.Upload(params)
		} else {
			msg, err = "", fmt.Errorf(
				"Exactly one of the source and destination url's must be a cloud url with the format cloud://...: %s, %s",
				src,
				dest,
			)
		}
		if err != nil {
			log.Fatal(msg, err)
		}
		writeOutput(msg)
	},
}

var LsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List objects",
	Long:  "List objects stored in a provider",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		p := mustGetStorageProvider(os.Getenv(storageProvEnv), os.Getenv(provEnv))
		url := convertURL(args[0], p.StorageURLPrefixReplacement())
		params := storage.NewLsParams(url)
		params.Recursive = lsRecursive
		msg, err := p.Ls(params)
		if err != nil {
			log.Fatal(msg, err)
		}
		writeOutput(msg)
	},
}

func mustGetStorageProvider (storageProvider string, provider string) storage.StorageProvider {
	p, err := getStorageProvider(storageProvider, provider)
	if err != nil {
		log.Fatal(err)
	}
	return p
}
