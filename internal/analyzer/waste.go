package analyzer

import (
	"fmt"
	"sort"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metricsv1beta1 "k8s.io/metrics/pkg/apis/metrics/v1beta1"
)

// PodWaste holds all display values and sort keys for one pod.
type PodWaste struct {
	Namespace    string `json:"namespace"`
	Pod          string `json:"pod"`
	CPURequest   string `json:"cpu_request"`
	CPUUsage     string `json:"cpu_usage"`
	CPUWaste     string `json:"cpu_waste"`
	CPURecommend string `json:"cpu_recommended"`
	MemRequest   string `json:"mem_request"`
	MemUsage     string `json:"mem_usage"`
	MemWaste     string `json:"mem_waste"`
	MemRecommend string `json:"mem_recommended"`

	// unexported — used only for sorting, not serialized
	cpuWasteRate float64
	memWasteRate float64
}

// Analyze matches pod specs with pod metrics and computes waste for every running pod.
func Analyze(pods []corev1.Pod, podMetrics map[string]metricsv1beta1.PodMetrics) []PodWaste {
	results := make([]PodWaste, 0, len(pods))

	for _, pod := range pods {
		if pod.Status.Phase != corev1.PodRunning {
			continue
		}

		key := pod.Namespace + "/" + pod.Name
		pm, hasMetrics := podMetrics[key]

		// Sum requests across all containers.
		var totalCPUReq, totalMemReq resource.Quantity
		for _, c := range pod.Spec.Containers {
			if v := c.Resources.Requests.Cpu(); v != nil {
				totalCPUReq.Add(*v)
			}
			if v := c.Resources.Requests.Memory(); v != nil {
				totalMemReq.Add(*v)
			}
		}

		// Sum usage across all containers (only when metrics are available).
		var totalCPUUse, totalMemUse resource.Quantity
		if hasMetrics {
			for _, c := range pm.Containers {
				if v := c.Usage.Cpu(); v != nil {
					totalCPUUse.Add(*v)
				}
				if v := c.Usage.Memory(); v != nil {
					totalMemUse.Add(*v)
				}
			}
		}

		w := PodWaste{Namespace: pod.Namespace, Pod: pod.Name}

		cpuReqM := totalCPUReq.MilliValue()
		cpuUseM := totalCPUUse.MilliValue()
		memReqB := totalMemReq.Value()
		memUseB := totalMemUse.Value()

		// CPU
		if cpuReqM == 0 {
			w.CPURequest = "none"
			w.CPUUsage = "-"
			w.CPUWaste = "n/a"
			w.CPURecommend = "n/a"
		} else {
			w.CPURequest = fmtMilliCPU(cpuReqM)
			if !hasMetrics {
				w.CPUUsage = "no metrics"
				w.CPUWaste = "n/a"
				w.CPURecommend = "n/a"
			} else {
				rate := wasteRate(cpuReqM, cpuUseM)
				w.cpuWasteRate = rate
				w.CPUUsage = fmtMilliCPU(cpuUseM)
				w.CPUWaste = fmtPercent(rate)
				w.CPURecommend = fmtMilliCPU(int64(float64(cpuUseM) * 1.5))
			}
		}

		// Memory
		if memReqB == 0 {
			w.MemRequest = "none"
			w.MemUsage = "-"
			w.MemWaste = "n/a"
			w.MemRecommend = "n/a"
		} else {
			w.MemRequest = fmtMiB(memReqB)
			if !hasMetrics {
				w.MemUsage = "no metrics"
				w.MemWaste = "n/a"
				w.MemRecommend = "n/a"
			} else {
				rate := wasteRate(memReqB, memUseB)
				w.memWasteRate = rate
				w.MemUsage = fmtMiB(memUseB)
				w.MemWaste = fmtPercent(rate)
				w.MemRecommend = fmtMiB(int64(float64(memUseB) * 1.5))
			}
		}

		results = append(results, w)
	}

	return results
}

// SortByWaste orders results by combined CPU+memory waste rate descending.
func SortByWaste(results []PodWaste) {
	sort.Slice(results, func(i, j int) bool {
		wi := results[i].cpuWasteRate + results[i].memWasteRate
		wj := results[j].cpuWasteRate + results[j].memWasteRate
		return wi > wj
	})
}

func wasteRate(request, usage int64) float64 {
	if request == 0 || usage >= request {
		return 0
	}
	return float64(request-usage) / float64(request) * 100
}

func fmtMilliCPU(m int64) string {
	return fmt.Sprintf("%dm", m)
}

func fmtMiB(b int64) string {
	return fmt.Sprintf("%dMi", b/(1024*1024))
}

func fmtPercent(p float64) string {
	return fmt.Sprintf("%.0f%%", p)
}
