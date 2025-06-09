package sinker

import (
	"context"
	"fmt"
	"time"

	"k8s.io/client-go/informers"
	eventsv1 "k8s.io/client-go/informers/events/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/AyCarlito/kube-event-sinker/pkg/logger"
	"github.com/AyCarlito/kube-event-sinker/pkg/sinker/sinks"
)

// resyncPeriod defines the time after which events in the local cache should be requeued.
// The duration here is around 290 years.
// It is intentionally high as the concept of requeing is made redundant by using the zeroOffsetSink.
const resyncPeriod time.Duration = 1<<63 - 1

// Sinker watches Kubernetes events and pushes them to a specified sink.
type Sinker struct {
	ctx       context.Context
	clientset *kubernetes.Clientset
	informer  eventsv1.EventInformer
}

// NewSinker returns a new *Sinker.
func NewSinker(ctx context.Context, kubeConfigPath, sinkName string) (*Sinker, error) {
	// Fetch in-cluster REST configuration. If this fails, use a local one in its place.
	restConfiguration, err := rest.InClusterConfig()
	if err != nil {
		kubeConfigLoadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
		kubeConfigLoadingRules.ExplicitPath = kubeConfigPath
		restConfiguration, err = clientcmd.BuildConfigFromFlags("", kubeConfigLoadingRules.GetDefaultFilename())
		if err != nil {
			return nil, fmt.Errorf("failed to get REST configuration: %v", err)
		}
	}

	clientset, err := kubernetes.NewForConfig(restConfiguration)
	if err != nil {
		return nil, fmt.Errorf("failed to create new clientset: %v", err)
	}

	// Select the Sink by name.
	// All Sinks provide a handler for the events.
	sink, err := sinks.NewSink(ctx, sinkName)
	if err != nil {
		return nil, err
	}

	// Wrap the sink with a zeroOffsetSink to filter out events that are assumed to have been previously handled.
	sink = sinks.NewSinkWithZeroOffset(sink)

	// Prefer use of informer factory to get a shared informer instead of getting an independant one.
	// Reduces memory footprint and number of connections to server.
	eventsInformer := informers.NewSharedInformerFactory(clientset, resyncPeriod).Events().V1().Events()

	// Add handlers for the events to the informer.
	eventsInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    sink.OnAdd,
		UpdateFunc: sink.OnUpdate,
		DeleteFunc: sink.OnDelete,
	})

	// We always setup a metrics sink.
	metricsSink := sinks.MetricsSink{}
	eventsInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    metricsSink.OnAdd,
		UpdateFunc: metricsSink.OnUpdate,
		DeleteFunc: metricsSink.OnDelete,
	})

	return &Sinker{
		ctx:       ctx,
		clientset: clientset,
		informer:  eventsInformer,
	}, nil
}

// Start starts the Sinker.
// The Sinker runs until its context is cancelled.
func (s *Sinker) Start() error {
	log := logger.LoggerFromContext(s.ctx)
	log.Info("Starting sinker")

	go s.informer.Informer().Run(s.ctx.Done())

	log.Info("Waiting for cache to sync")
	if !cache.WaitForCacheSync(s.ctx.Done(), s.informer.Informer().HasSynced) {
		return fmt.Errorf("failed to wait for cache to sync")
	}
	log.Info("Cache synced")
	<-s.ctx.Done()
	s.shutdown()
	return nil
}

// shutdown shuts down the Sinker.
// It blocks until the underlying Informer has stopped.
func (s *Sinker) shutdown() {
	log := logger.LoggerFromContext(s.ctx)
	log.Info("Shutting down sinker")
	for {
		if s.informer.Informer().IsStopped() {
			return
		}
	}
}
