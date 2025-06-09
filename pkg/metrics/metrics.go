package metrics

import "github.com/prometheus/client_golang/prometheus"

func init() {
	prometheus.MustRegister(KubernetesEvents)
}

// KubernetesEvents is a prometheus CounterVec that totals the Kubernetes events handled.
var KubernetesEvents = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "kube_event_sinker_events_total",
		Help: "Total of Kubernetes events handled.",
	},
	[]string{
		"regarding_kind",
		"regarding_name",
		"regarding_namespace",
		"reason",
		"type",
	},
)
