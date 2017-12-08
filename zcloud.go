package main

import (
	"github.com/nathanwilk7/zcloud/cmd"
)

func main () {
	rootCmd := cmd.RootCmd
	rootCmd.AddCommand(cmd.StorageCmd)
	rootCmd.Execute()
}
