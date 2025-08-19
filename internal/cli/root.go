package cli

import (
	"fmt"
	"os"
	
	"github.com/spf13/cobra"
	"github.com/autodevops/verifier-go/internal/util"
	"github.com/sirupsen/logrus"
)

var verbose bool

var rootCmd = &cobra.Command{
	Use:   "verifier",
	Short: "AI-powered code verification CLI (Go Version)",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if verbose {
			util.Log.SetLevel(logrus.DebugLevel)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		util.Log.WithError(err).Fatal("CLI execution failed")
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")
}
