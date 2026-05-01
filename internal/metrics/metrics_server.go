package metrics

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	metricsv1beta1 "k8s.io/metrics/pkg/apis/metrics/v1beta1"

	"github.com/k8s-waste-killer/internal/k8s"
)

// FetchPodMetrics returns a map keyed by "namespace/podName" for fast lookup.
func FetchPodMetrics(client *k8s.Client, namespace string) (map[string]metricsv1beta1.PodMetrics, error) {
	list, err := client.Metrics.MetricsV1beta1().PodMetricses(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	result := make(map[string]metricsv1beta1.PodMetrics, len(list.Items))
	for _, m := range list.Items {
		key := m.Namespace + "/" + m.Name
		result[key] = m
	}
	return result, nil
}
