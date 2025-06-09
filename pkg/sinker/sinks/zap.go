package sinks

import (
	"context"

	"github.com/AyCarlito/kube-event-sinker/pkg/logger"
	"go.uber.org/zap"
	eventsv1 "k8s.io/api/events/v1"
)

// zapSink is a sink that records events through a zap logger.
type zapSink struct {
	ctx context.Context
}

// OnAdd handles Add events.
func (z *zapSink) OnAdd(obj interface{}) {
	z.handle(obj)
}

// OnUpdate handles Update events.
func (z *zapSink) OnUpdate(oldObj, newObj interface{}) {
	z.handle(newObj)
}

// OnDelete handles Delete events.
func (z *zapSink) OnDelete(obj interface{}) {
	z.handle(obj)
}

// handle handles an event.
func (z *zapSink) handle(obj interface{}) {
	event := obj.(*eventsv1.Event)
	log := logger.LoggerFromContext(z.ctx).With(
		zap.String("kind", event.Regarding.Kind),
		zap.String("name", event.Regarding.Name),
		zap.String("namespace", event.Regarding.Namespace),
		zap.String("reason", event.Reason),
		zap.String("type", event.Type),
	)
	log.Info("Handling")
}
