package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/autodevops/verifier-go/internal/config"
	"github.com/spf13/cobra"
)

var force bool

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize verifier in the current repository",
	RunE: func(cmd *cobra.Command, args []string) error {
		configPath := filepath.Join(".verifier", "config.yaml")
		if _, err := os.Stat(configPath); err == nil && !force {
			fmt.Println("Verifier already initialized. Use --force to overwrite.")
			return nil
		}

		fmt.Println("Initializing verifier...")

		defaultConfig := &config.Config{
			Models: config.Models{
				Primary:  "claude-3-5-sonnet-20240620",
				Fallback: "claude-3-haiku-20240307",
			},
			Providers: config.Providers{}, // API keys should be in .env
			Budgets: config.Budgets{
				DailyTokens:     100000,
				PerCommitTokens: 5000,
				MonthlyCost:     100,
			},
			Thresholds: config.Thresholds{
				DriftScore:    30,
				SecurityRisk:  5,
				CoverageDelta: -5,
			},
			Hooks: map[string][]string{
				"pre-commit": {"lint", "security-scan"},
			},
		}

		if err := config.Save(defaultConfig); err != nil {
			return fmt.Errorf("failed to save config: %w", err)
		}

		envPath := filepath.Join(".verifier", ".env")
		envContent := "ANTHROPIC_API_KEY=\"YOUR_API_KEY_HERE\"\n"
		if err := os.WriteFile(envPath, []byte(envContent), 0600); err != nil {
			return fmt.Errorf("failed to create .env file: %w", err)
		}
		
		fmt.Println("✓ Verifier initialized successfully!")
		return nil
	},
}

func init() {
	initCmd.Flags().BoolVarP(&force, "force", "f", false, "Overwrite existing configuration")
	rootCmd.AddCommand(initCmd)
}
