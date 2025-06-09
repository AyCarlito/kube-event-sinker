package sinks

// NullSink is a sink that does nothing.
type NullSink struct{}

// OnAdd does nothing with an Add event.
func (n *NullSink) OnAdd(obj interface{}) {}

// OnUpdate does nothing with an Update event.
func (n *NullSink) OnUpdate(oldObj, newObj interface{}) {}

// OnDelete does nothing with a Delete event.
func (n *NullSink) OnDelete(obj interface{}) {}
