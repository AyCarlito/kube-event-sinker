package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"

	"github.com/AyCarlito/kube-event-sinker/pkg/logger"
	"github.com/AyCarlito/kube-event-sinker/pkg/sinker"
)

func init() {
	rootCmd.PersistentFlags().StringVar(&metricsAddress, "metrics-bind-address", ":9111", "The address the metric endpoint binds to.")
	rootCmd.PersistentFlags().StringVar(&kubeConfigPath, "kubeconfig", "", "Path to a kubeconfig file.")
	rootCmd.PersistentFlags().StringVar(&sinkName, "sink", "null", "Sink that events should be pushed to.")
}

// CLI Flags
var (
	metricsAddress string
	kubeConfigPath string
	sinkName       string
)

var rootCmd = &cobra.Command{
	Use:           "kube-event-sinker",
	Short:         "Watch and push Kubernetes events to a sink.",
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Build logger.
		log, err := logger.NewZapConfig().Build()
		if err != nil {
			panic(fmt.Errorf("failed to build zap logger: %v", err))
		}
		cmd.SetContext(logger.ContextWithLogger(cmd.Context(), log))

		sinker, err := sinker.NewSinker(cmd.Context(), kubeConfigPath, sinkName)
		if err != nil {
			panic(fmt.Errorf("failed to create new sinker: %v", err))
		}

		go func() {
			// Metrics.
			mux := http.NewServeMux()
			mux.Handle("/metrics", promhttp.Handler())
			err := http.ListenAndServe(metricsAddress, mux)
			if err != nil {
				panic(fmt.Errorf("failed to start metrics server: %v", err))
			}
		}()

		return sinker.Start()
	},
}

func Execute() {
	ctx, cxl := context.WithCancel(context.Background())
	defer cxl()

	// Following a signal, cancel the context to initiate graceful shutdown.
	// This produces one of two desired behaviours:
	// 1) Exit after shutdown is complete.
	// 2) Exit on a subsequent signal.
	sigChan := make(chan os.Signal, 2)
	signal.Notify(sigChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sigChan
		cxl()
		<-sigChan
		os.Exit(-1)
	}()

	err := rootCmd.ExecuteContext(ctx)
	if err != nil {
		// By default, cobra prints the error and usage string on every error.
		// We only desire this behaviour in the case where command line parsing fails e.g. unknown command or flag.
		// Cobra does not provide a mechanism for achieving this fine grain control, so we implement our own.
		if strings.Contains(err.Error(), "command") || strings.Contains(err.Error(), "flag") {
			// Parsing errors are printed along with the usage string.
			fmt.Println(err.Error())
			fmt.Println(rootCmd.UsageString())
		} else {
			// Other errors logged, no usage string displayed.
			log := logger.LoggerFromContext(rootCmd.Context())
			log.Error(err.Error())
		}
		os.Exit(1)
	}
}
