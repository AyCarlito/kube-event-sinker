package sinks

import (
	"github.com/prometheus/client_golang/prometheus"
	eventsv1 "k8s.io/api/events/v1"

	"github.com/AyCarlito/kube-event-sinker/pkg/metrics"
)

// MetricsSink is a sink that generates prometheus metrics.
type MetricsSink struct{}

// OnAdd handles Add events.
func (m *MetricsSink) OnAdd(obj interface{}) {
	m.generate(obj)
}

// OnUpdate handles Update events.
func (m *MetricsSink) OnUpdate(oldObj, newObj interface{}) {
	m.generate(newObj)
}

// OnDelete handles Delete events.
func (m *MetricsSink) OnDelete(obj interface{}) {
	m.generate(obj)
}

// generate generates prometheus metrics.
func (m *MetricsSink) generate(obj interface{}) {
	event := obj.(*eventsv1.Event)
	metrics.KubernetesEvents.With(prometheus.Labels{
		"regarding_kind":      event.Regarding.Kind,
		"regarding_name":      event.Regarding.Name,
		"regarding_namespace": event.Regarding.Namespace,
		"reason":              event.Reason,
		"type":                event.Type,
	}).Inc()
}
