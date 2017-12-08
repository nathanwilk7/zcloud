package cmd

import (
	"log"
	
	"github.com/spf13/cobra"

	"github.com/nathanwilk7/zcloud/providers"
	"github.com/nathanwilk7/zcloud/storage"
)

var StorageCmd = &cobra.Command{
	Use:   "storage",
	Short: "Store stuff",
	Long:  "Store stuff long description",
}

func init () {
	CpCmd.Flags().BoolVarP(&cpRecursive, "recursive", "r", false, "Recursively copy from src")
	StorageCmd.AddCommand(CpCmd)
}

var cpRecursive bool

var CpCmd = &cobra.Command{
	Use:   "cp",
	Short: "Copy stuff",
	Long:  `Copy stuff long description`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		src := args[0]
		dest := args[1]
		params := storage.NewCpParams(src, dest)
		params.Recursive = cpRecursive
		msg, err := providers.ProviderInstance.Cp(params)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(msg)
	},
}
