package cmd

import (
	"fmt"
	
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "zcloud",
	Short: "zCloud makes using the cloud better",
	Long: "zCloud provides a layer of abstraction between clients and cloud providers",
}

var quiet bool

func init () {
	RootCmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false, "Don't print output")
}

func writeOutput (msg string) {
	if !quiet {
		fmt.Println(msg)
	}
}
