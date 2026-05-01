package output

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/k8s-waste-killer/internal/analyzer"
)

func PrintTable(results []analyzer.PodWaste) {
	if len(results) == 0 {
		fmt.Println("No running pods found.")
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "NAMESPACE\tPOD\tCPU_REQ\tCPU_USE\tCPU_WASTE\tCPU_RECOMMEND\tMEM_REQ\tMEM_USE\tMEM_WASTE\tMEM_RECOMMEND")

	for _, r := range results {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
			r.Namespace, r.Pod,
			r.CPURequest, r.CPUUsage, r.CPUWaste, r.CPURecommend,
			r.MemRequest, r.MemUsage, r.MemWaste, r.MemRecommend,
		)
	}

	w.Flush()
}
