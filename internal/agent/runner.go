package agent

import (
	"fmt"
	"time"
	
	"github.com/autodevops/verifier-go/internal/config"
	"github.com/autodevops/verifier-go/internal/storage"
)

type AgentRunner struct {
	cfg    *config.Config
	metrics *storage.MetricsStore
}

func NewAgentRunner(cfg *config.Config) *AgentRunner {
	return &AgentRunner{
		cfg:    cfg,
		metrics: storage.NewMetricsStore(),
	}
}

func (r *AgentRunner) RunAgent(id string, ctx AgentContext) (*AgentResult, error) {
	// Check budget before running
	todaysMetrics, _ := r.metrics.GetMetrics(24 * time.Hour)
	tokensUsedToday := 0
	for _, m := range todaysMetrics {
		tokensUsedToday += m.TokensUsed
	}
	if tokensUsedToday >= r.cfg.Budgets.DailyTokens {
		return &AgentResult{
			AgentID:   id,
			Status:    "skipped",
			Error:     "Daily token budget exhausted",
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		}, nil
	}

	agent, err := GetAgent(id, r.cfg)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	result, err := agent.Execute(ctx)
	duration := time.Since(start)

	if err != nil {
		result = &AgentResult{
			AgentID:   id,
			Status:    "failure",
			Error:     err.Error(),
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		}
	}

	// Record metrics
	_ = r.metrics.Record(storage.Metric{
		AgentID:    id,
		Timestamp:  time.Now(),
		TokensUsed: result.TokensUsed,
		Cost:       result.Cost,
		Result:     result.Status,
		DurationMs: duration.Milliseconds(),
	})

	return result, nil
}
