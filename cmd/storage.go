package cmd

import (
	"log"

	"github.com/nathanwilk7/zcloud/providers"
	"github.com/nathanwilk7/zcloud/storage"

	"github.com/spf13/cobra"
)

var StorageCmd = &cobra.Command{
	Use:   "storage",
	Short: "Store objects",
	Long:  "Store objects long description",
}

func init () {
	CpCmd.Flags().BoolVarP(&recursive, "recursive", "r", false, "Recursively copy from src")
	StorageCmd.AddCommand(CpCmd)

	LsCmd.Flags().BoolVarP(&recursive, "recursive", "r", false, "Recursively list")
	StorageCmd.AddCommand(LsCmd)
}

var recursive bool

var CpCmd = &cobra.Command{
	Use:   "cp",
	Short: "Copy objects",
	Long:  `Copy objects long description`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		src, dest := parseSrcDestUrl(args)
		params := storage.NewCpParams(src, dest)
		params.Recursive = recursive
		p, err := providers.GetStorageProvider()
		if err != nil {
			log.Fatal(err)
		}
		msg, err := p.Cp(params)
		if err != nil {
			log.Fatal(msg, err)
		}
		logOutput(quiet, msg)
	},
}

var LsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List objects",
	Long:  `List objects long description`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		url := parseUrl(args)
		params := storage.NewLsParams(url)
		params.Recursive = recursive
		p, err := providers.GetStorageProvider()
		if err != nil {
			log.Fatal(err)
		}
		msg, err := p.Ls(params)
		if err != nil {
			log.Fatal(msg, err)
		}
		logOutput(quiet, msg)
	},
}
