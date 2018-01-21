package cmd

import (
	"log"
	
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "zcloud",
	Short: "zCloud makes using the cloud better",
	Long: "zCloud provides a layer of abstraction between clients using cloud providers to prevent cloud provider coupling",
}

func init () {
	CpCmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "Don't print output")
}

var quiet bool

func logOutput (q bool, msg string) {
	if !q {
		log.Println(msg)
	}
}

func parseSrcDestUrl (args []string) (string, string) {
	src := args[0]
	dest := args[1]
	return src, dest
}

func parseUrl (args []string) string {
	return args[0]
}
