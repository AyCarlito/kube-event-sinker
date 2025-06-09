package sinks

import (
	"context"
	"fmt"
)

const (
	nullSinkName string = "null"
	zapSinkName  string = "zap"
)

// Sink handles events.
// Handlers are informational only, they must not modify the event object.
type Sink interface {
	OnAdd(obj interface{})
	OnUpdate(oldObj, newObj interface{})
	OnDelete(obj interface{})
}

// NewSink returns the Sink corresponding to the provided sink name.
func NewSink(ctx context.Context, name string) (Sink, error) {
	// TODO: Should the metrics sink be a wrapper?
	// TODO: Should also skip events from before the application started.
	switch name {
	case nullSinkName:
		return &nullSink{}, nil
	case zapSinkName:
		return &zapSink{ctx: ctx}, nil
	default:
		return nil, fmt.Errorf("unrecognised sink: %s", name)
	}
}
