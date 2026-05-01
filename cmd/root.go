package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "k8s-waste-killer",
	Short: "Analyze Kubernetes resource waste",
	Long: `k8s-waste-killer compares CPU and memory requests to actual usage
and shows which pods are over-provisioned, along with right-sized recommendations.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(scanCmd)
}
