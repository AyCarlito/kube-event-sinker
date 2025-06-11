package sinks

import (
	"github.com/prometheus/client_golang/prometheus"
	eventsv1 "k8s.io/api/events/v1"

	"github.com/AyCarlito/kube-event-sinker/pkg/metrics"
)

// metricsSink is a sink that generates prometheus metrics.
// It acts as a wrapper around a provided Sink.
type metricsSink struct {
	sink Sink
}

// OnAdd handles Add events.
func (m *metricsSink) OnAdd(obj interface{}) {
	m.handle(obj)
	m.sink.OnAdd(obj)
}

// OnUpdate handles Update events.
func (m *metricsSink) OnUpdate(oldObj, newObj interface{}) {
	m.handle(newObj)
	m.sink.OnUpdate(oldObj, newObj)
}

// OnDelete handles Delete events.
func (m *metricsSink) OnDelete(obj interface{}) {
	m.handle(obj)
	m.sink.OnDelete(obj)
}

// handle handles an event.
func (m *metricsSink) handle(obj interface{}) {
	event := obj.(*eventsv1.Event)
	metrics.KubernetesEvents.With(prometheus.Labels{
		"kind":   event.Regarding.Kind,
		"reason": event.Reason,
		"type":   event.Type,
	}).Inc()
}

// NewSinkWithMetrics wraps a provided sink a *metricsSink.
func NewSinkWithMetrics(sink Sink) Sink {
	return &metricsSink{sink: sink}
}
