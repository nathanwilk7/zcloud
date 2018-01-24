package cmd

import (
	"fmt"
	"strings"

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

const cloudStr = "cloud"
const cloudURLPrefix = cloudStr + "://"

func convertURL (url, replacement string) string {
	if isCloudURL(url) {
		return strings.Replace(url, cloudStr, replacement, 1)
	}
	return url
}

func isCloudURL (url string) bool {
	if len(url) > len(cloudURLPrefix) {
		if url[:len(cloudURLPrefix)] == cloudURLPrefix {
			return true
		}
	}
	return false
}
