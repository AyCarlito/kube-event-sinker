package sinker

import (
	"context"
	"fmt"
	"time"

	"github.com/AyCarlito/kube-event-sinker/pkg/logger"
	"k8s.io/client-go/informers"
	eventsv1 "k8s.io/client-go/informers/events/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

// Sinker watches Kubernetes events and pushes them to a specified sink.
type Sinker struct {
	ctx       context.Context
	clientset *kubernetes.Clientset
	informer  eventsv1.EventInformer
}

// NewSinker returns a new *Sinker.
func NewSinker(ctx context.Context, kubeConfigPath string) (*Sinker, error) {
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

	// Prefer use of informer factory to get a shared informer instead of getting an independant one.
	// Reduces memory footprint and number of connections to server.
	// TODO: Resync duration should come from CLI flag.
	eventsInformer := informers.NewSharedInformerFactory(clientset, 1*time.Hour).Events().V1().Events()

	// Add handlers for the events to the informer.
	// Handlers are informational only, they must not modify the event object.
	eventsInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		// TODO: Placeholders for handlers from a sink.
		// The sinks themselves shouldn't do metrics. The sink handlers should be wrapped.
		// Should also skip events from before the application started.
		AddFunc: func(obj interface{}) {
			fmt.Println(obj)
			fmt.Println("ADD")
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			fmt.Println("UPDATE")
		},
		DeleteFunc: func(obj interface{}) {
			fmt.Println("DELETE")
		},
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
	return nil
}
