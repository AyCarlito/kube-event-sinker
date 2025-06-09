package sinks

import (
	"time"

	eventsv1 "k8s.io/api/events/v1"
)

// zeroOffsetSink is a sink that filters out events that are assumed to have already been handled.
// It acts as a wrapper around a provided Sink.
type zeroOffsetSink struct {
	sink Sink
	t    time.Time
}

// OnAdd handles Add events.
func (z *zeroOffsetSink) OnAdd(obj interface{}) {
	event := obj.(*eventsv1.Event)
	// Ignore events with a last timestamp earlier than the time the zeroOffsetSink was created.
	if z.t.After(event.DeprecatedLastTimestamp.Time) {
		return
	}
	z.sink.OnAdd(obj)
}

// OnUpdate handles Update events.
func (z *zeroOffsetSink) OnUpdate(oldObj, newObj interface{}) {
	event := newObj.(*eventsv1.Event)
	// Ignore events with a last timestamp earlier than the time the zeroOffsetSink was created.
	if z.t.After(event.DeprecatedLastTimestamp.Time) {
		return
	}
	z.sink.OnUpdate(oldObj, newObj)
}

// OnDelete handles Delete events.
func (z *zeroOffsetSink) OnDelete(obj interface{}) {
	event := obj.(*eventsv1.Event)
	// Ignore events with a last timestamp earlier than the time the zeroOffsetSink was created.
	if z.t.After(event.DeprecatedLastTimestamp.Time) {
		return
	}
	z.sink.OnDelete(obj)
}

// NewSinkWithZeroOffset wraps a provided sink a *zeroOffsetSink.
func NewSinkWithZeroOffset(sink Sink) Sink {
	return &zeroOffsetSink{
		sink: sink,
		t:    time.Now().UTC(),
	}
}
