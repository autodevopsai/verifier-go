package cli

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/autodevopsai/verifier-go/internal/storage"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var period string
var format string

var tokenUsageCmd = &cobra.Command{
	Use:   "token-usage",
	Short: "Show token usage statistics",
	RunE: func(cmd *cobra.Command, args []string) error {
		var duration time.Duration
		switch period {
		case "hourly":
			duration = time.Hour
		case "daily":
			duration = 24 * time.Hour
		case "weekly":
			duration = 7 * 24 * time.Hour
		case "monthly":
			duration = 30 * 24 * time.Hour
		default:
			return fmt.Errorf("invalid period: %s", period)
		}

		store := storage.NewMetricsStore()
		metrics, err := store.GetMetrics(duration)
		if err != nil {
			return err
		}

		usage := make(map[string]map[string]float64)
		totalTokens := 0
		totalCost := 0.0

		for _, m := range metrics {
			if _, ok := usage[m.AgentID]; !ok {
				usage[m.AgentID] = make(map[string]float64)
			}
			usage[m.AgentID]["tokens"] += float64(m.TokensUsed)
			usage[m.AgentID]["cost"] += m.Cost
			usage[m.AgentID]["calls"]++
			totalTokens += m.TokensUsed
			totalCost += m.Cost
		}

		// Table output
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Agent", "Calls", "Tokens", "Cost"})
		for agentID, data := range usage {
			row := []string{
				agentID,
				strconv.Itoa(int(data["calls"])),
				strconv.Itoa(int(data["tokens"])),
				fmt.Sprintf("$%.4f", data["cost"]),
			}
			table.Append(row)
		}
		table.Render()

		fmt.Printf("\nTotal Tokens: %d\n", totalTokens)
		fmt.Printf("Total Cost: $%.4f\n", totalCost)
		
		return nil
	},
}

func init() {
	tokenUsageCmd.Flags().StringVarP(&period, "period", "p", "daily", "Time period (hourly|daily|weekly|monthly)")
	tokenUsageCmd.Flags().StringVarP(&format, "format", "f", "table", "Output format (table|json)")
	rootCmd.AddCommand(tokenUsageCmd)
}
