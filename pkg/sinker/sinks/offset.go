package sinks

import (
	"time"

	eventsv1 "k8s.io/api/events/v1"
)

// offsetSink is a sink that filters out events that are assumed to have already been handled.
// It acts as a wrapper around a provided Sink.
type offsetSink struct {
	sink Sink
	t    time.Time
}

// OnAdd handles Add events.
func (o *offsetSink) OnAdd(obj interface{}) {
	event := obj.(*eventsv1.Event)
	// Ignore events with a last timestamp earlier than the time the offsetSink was created.
	if o.t.After(event.DeprecatedLastTimestamp.Time) {
		return
	}
	o.sink.OnAdd(obj)
}

// OnUpdate handles Update events.
func (o *offsetSink) OnUpdate(oldObj, newObj interface{}) {
	event := newObj.(*eventsv1.Event)
	// Ignore events with a last timestamp earlier than the time the offsetSink was created.
	if o.t.After(event.DeprecatedLastTimestamp.Time) {
		return
	}
	o.sink.OnUpdate(oldObj, newObj)
}

// OnDelete handles Delete events.
func (o *offsetSink) OnDelete(obj interface{}) {
	event := obj.(*eventsv1.Event)
	// Ignore events with a last timestamp earlier than the time the offsetSink was created.
	if o.t.After(event.DeprecatedLastTimestamp.Time) {
		return
	}
	o.sink.OnDelete(obj)
}

// NewSinkWithOffset wraps a provided sink an *offsetSink.
func NewSinkWithOffset(sink Sink) Sink {
	return &offsetSink{
		sink: sink,
		t:    time.Now().UTC(),
	}
}
