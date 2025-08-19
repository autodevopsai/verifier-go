package cli

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/autodevops/verifier-go/internal/agent"
	"github.com/autodevops/verifier-go/internal/config"
	"github.com/autodevops/verifier-go/internal/context"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run [agent-id]",
	Short: "Run a specified verifier agent",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		agentID := args[0]
		
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load configuration: %w. Please run 'verifier init'", err)
		}

		if _, err := agent.GetAgent(agentID, cfg); err != nil {
			available := strings.Join(agent.ListAgents(), ", ")
			return fmt.Errorf("%w. Available agents: %s", err, available)
		}

		fmt.Printf("Running agent: %s...\n", agentID)

		ctx, err := context.CollectGitContext()
		if err != nil {
			fmt.Printf("Warning: could not collect git context: %v\n", err)
		}

		runner := agent.NewAgentRunner(cfg)
		result, err := runner.RunAgent(agentID, ctx)
		if err != nil {
			return fmt.Errorf("agent execution failed: %w", err)
		}

		output, _ := json.MarshalIndent(result, "", "  ")
		fmt.Println(string(output))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
