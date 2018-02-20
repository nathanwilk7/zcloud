package cmd

import (
	"os"

	"github.com/nathanwilk7/zcloud/controller"
	"github.com/nathanwilk7/zcloud/out"
	
	"github.com/spf13/cobra"
)

var StorageCmd = &cobra.Command{
	Use:   "storage",
	Short: "Blob storage",
	Long:  "Perform various operations related to blob storage",
}

var cpRecursive bool
var lsRecursive bool
var rmRecursive bool

func init () {
	CpCmd.Flags().BoolVarP(&cpRecursive, "recursive", "r", false, "Recursively copy from src")
	StorageCmd.AddCommand(CpCmd)
	
	LsCmd.Flags().BoolVarP(&lsRecursive, "recursive", "r", false, "Recursively list")
	StorageCmd.AddCommand(LsCmd)

	RmCmd.Flags().BoolVarP(&rmRecursive, "recursive", "r", false, "Recursively remove")
	StorageCmd.AddCommand(RmCmd)
	
	RootCmd.AddCommand(StorageCmd)
}

const storageProvEnv = "ZCLOUD_STORAGE_PROV"
var storageProv = os.Getenv(storageProvEnv)

func getProvParamsFromEnv () controller.ProvParams {
	return controller.ProvParams{
		Name: getStorageProv(storageProv, prov),
		AwsId: awsId,
		AwsSecret: awsSecret,
		AwsRegion: awsRegion,
	}
}

func getStorageProv (storageProv, prov string) string {
	if storageProv != "" {
		return storageProv
	}
	return prov
}

var CpCmd = &cobra.Command{
	Use:   "cp",
	Short: "Copy objects",
	Long:  "Copy objects to/from a provider",
	Args: cobra.ExactArgs(2),
	Run: func (cmd *cobra.Command, args []string) {
		pp := getProvParamsFromEnv()
		cp := controller.CpParams{
			Src: args[0],
			Dest: args[1],
			Recursive: cpRecursive,
		}
		controller.Cp(pp, cp, out.New())
	},
}

var LsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List objects",
	Long:  "List objects stored in a provider",
	Args: cobra.ExactArgs(1),
	Run: func (cmd *cobra.Command, args []string) {
		pp := getProvParamsFromEnv()
		lp := controller.LsParams{
			Url: args[0],
			Recursive: lsRecursive,
		}
		controller.Ls(pp, lp, out.New())
	},
}

var RmCmd = &cobra.Command{
	Use:   "rm",
	Short: "Remove objects",
	Long:  "Remove objects stored in a provider",
	Args: cobra.ExactArgs(1),
	Run: func (cmd *cobra.Command, args []string) {
		pp := getProvParamsFromEnv()
		rp := controller.RmParams{
			Url: args[0],
			Recursive: rmRecursive,
		}
		controller.Rm(pp, rp, out.New())
	},
}
