package sinks

// nullSink is a sink that does nothing.
type nullSink struct{}

// OnAdd handles Add events.
func (n *nullSink) OnAdd(obj interface{}) {}

// OnUpdate handles Update events.
func (n *nullSink) OnUpdate(oldObj, newObj interface{}) {}

// OnDelete handles Delete events.
func (n *nullSink) OnDelete(obj interface{}) {}
