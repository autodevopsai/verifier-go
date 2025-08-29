package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/autodevopsai/verifier-go/internal/config"
	"github.com/spf13/cobra"
)

type Check struct {
	Name    string
	OK      bool
	Message string
}

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Verify environment and configuration",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🩺 Running system diagnostics...")
		var checks []Check
		failures := 0

		// Go version
		checks = append(checks, Check{"Go version", true, runtime.Version()})

		// Config directory
		dirExists := checkPathExists(".verifier")
		checks = append(checks, Check{"Config directory .verifier/", dirExists, ""})

		// Config file
		configExists := checkPathExists(filepath.Join(".verifier", "config.yaml"))
		checks = append(checks, Check{"Config file present", configExists, ""})

		// Config parsing and API key check
		cfg, err := config.Load()
		checks = append(checks, Check{"Config parsing", err == nil, ""})
		if err == nil {
			keyOk := cfg.Providers.OpenAI.APIKey != "" || cfg.Providers.Anthropic.APIKey != ""
			checks = append(checks, Check{"Provider API key", keyOk, ""})
		}

		// Git binary
		_, err = exec.LookPath("git")
		checks = append(checks, Check{"git binary in PATH", err == nil, ""})

		for _, c := range checks {
			icon := "✅"
			if !c.OK {
				icon = "❌"
				failures++
			}
			fmt.Printf("%s %-30s %s\n", icon, c.Name, c.Message)
		}

		if failures > 0 {
			fmt.Printf("\nFound %d issue(s). Please review the output above.\n", failures)
		} else {
			fmt.Println("\nAll checks passed. Verifier is ready!")
		}
	},
}

func checkPathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func init() {
	rootCmd.AddCommand(doctorCmd)
}
