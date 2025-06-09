package sinks

import "fmt"

const (
	nullSink string = "null"
)

// Sink handles events.
// Handlers are informational only, they must not modify the event object.
type Sink interface {
	OnAdd(obj interface{})
	OnUpdate(oldObj, newObj interface{})
	OnDelete(obj interface{})
}

// NewSink returns the Sink corresponding to the provided sink name.
func NewSink(name string) (Sink, error) {
	// TODO: Metrics. The sink handlers should be wrapped.
	// TODO: Should also skip events from before the application started.
	switch name {
	case nullSink:
		return &NullSink{}, nil
	default:
		return nil, fmt.Errorf("unrecognised sink: %s", name)
	}
}
