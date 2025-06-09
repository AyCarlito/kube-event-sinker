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
	m.handle(obj)
}

// OnUpdate handles Update events.
func (m *MetricsSink) OnUpdate(oldObj, newObj interface{}) {
	m.handle(newObj)
}

// OnDelete handles Delete events.
func (m *MetricsSink) OnDelete(obj interface{}) {
	m.handle(obj)
}

// handle handles an event.
func (m *MetricsSink) handle(obj interface{}) {
	event := obj.(*eventsv1.Event)
	metrics.KubernetesEvents.With(prometheus.Labels{
		"kind":      event.Regarding.Kind,
		"name":      event.Regarding.Name,
		"namespace": event.Regarding.Namespace,
		"reason":    event.Reason,
		"type":      event.Type,
	}).Inc()
}
