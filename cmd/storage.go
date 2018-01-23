package cmd

import (
	"log"

	"github.com/nathanwilk7/zcloud/providers"
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
}

var CpCmd = &cobra.Command{
	Use:   "cp",
	Short: "Copy objects",
	Long:  "Copy objects to/from a provider",
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		src, dest := args[0], args[1]
		params := storage.NewCpParams(src, dest)
		params.Recursive = cpRecursive
		p := mustGetStorageProvider()
		msg, err := p.Cp(params)
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
		url := args[0]
		params := storage.NewLsParams(url)
		params.Recursive = lsRecursive
		p := mustGetStorageProvider()
		msg, err := p.Ls(params)
		if err != nil {
			log.Fatal(msg, err)
		}
		writeOutput(msg)
	},
}

func mustGetStorageProvider () storage.StorageProvider {
	p, err := providers.GetStorageProvider()
	if err != nil {
		log.Fatal(err)
	}
	return p
}
