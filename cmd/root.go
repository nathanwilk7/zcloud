package cmd

import (
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "zcloud",
	Short: "zCloud makes using the cloud better",
	Long: "zCloud provides a layer of abstraction between clients using cloud providers to prevent cloud provider coupling",
}
