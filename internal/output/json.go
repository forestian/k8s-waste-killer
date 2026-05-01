package output

import (
	"encoding/json"
	"os"

	"github.com/k8s-waste-killer/internal/analyzer"
)

func PrintJSON(results []analyzer.PodWaste) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(results)
}
