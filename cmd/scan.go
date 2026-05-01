package cmd

import (
	"fmt"
	"os"

	"github.com/k8s-waste-killer/internal/analyzer"
	"github.com/k8s-waste-killer/internal/k8s"
	"github.com/k8s-waste-killer/internal/metrics"
	"github.com/k8s-waste-killer/internal/output"
	"github.com/spf13/cobra"
)

var (
	namespace  string
	topN       int
	outputFmt  string
	kubeconfig string
)

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan pods for resource waste",
	Long:  `Compares CPU and memory requests against actual usage from metrics-server and ranks pods by waste.`,
	RunE:  runScan,
}

func init() {
	scanCmd.Flags().StringVarP(&namespace, "namespace", "n", "", "Kubernetes namespace to scan (default: all namespaces)")
	scanCmd.Flags().IntVar(&topN, "top", 0, "Show only the top N pods by waste (0 = show all)")
	scanCmd.Flags().StringVarP(&outputFmt, "output", "o", "table", "Output format: table or json")
	scanCmd.Flags().StringVar(&kubeconfig, "kubeconfig", "", "Path to kubeconfig file (defaults to $KUBECONFIG or ~/.kube/config)")
}

func runScan(cmd *cobra.Command, args []string) error {
	client, err := k8s.NewClient(kubeconfig)
	if err != nil {
		return fmt.Errorf("failed to create Kubernetes client: %w", err)
	}

	pods, err := k8s.ListPods(client, namespace)
	if err != nil {
		return fmt.Errorf("failed to list pods: %w", err)
	}

	podMetrics, err := metrics.FetchPodMetrics(client, namespace)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: could not fetch metrics (is metrics-server installed?): %v\n", err)
		podMetrics = nil
	}

	results := analyzer.Analyze(pods, podMetrics)
	analyzer.SortByWaste(results)

	if topN > 0 && topN < len(results) {
		results = results[:topN]
	}

	switch outputFmt {
	case "json":
		return output.PrintJSON(results)
	default:
		output.PrintTable(results)
		return nil
	}
}
