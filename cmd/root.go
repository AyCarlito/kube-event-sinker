package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/AyCarlito/kube-event-sinker/pkg/logger"
	"github.com/AyCarlito/kube-event-sinker/pkg/sinker"
)

func init() {
	rootCmd.PersistentFlags().StringVar(&kubeConfigPath, "kubeconfig", "", "Path to a kubeconfig file.")
}

// CLI Flags
var (
	kubeConfigPath string
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

		sinker, err := sinker.NewSinker(cmd.Context(), kubeConfigPath)
		if err != nil {
			panic(fmt.Errorf("failed to create new sinker: %v", err))
		}

		return sinker.Start()
	},
}

func Execute() {
	err := rootCmd.Execute()
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
