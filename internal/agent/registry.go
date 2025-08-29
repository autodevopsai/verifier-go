package agent

import (
	"fmt"
	
	"github.com/autodevopsai/verifier-go/internal/config"
)

var agentInitializers map[string]func(cfg *config.Config) Agent

func init() {
	agentInitializers = make(map[string]func(cfg *config.Config) Agent)
	Register("lint", func(_ *config.Config) Agent { return NewLintAgent() })
	Register("security-scan", NewSecurityScanAgent)
}

// Register adds a new agent initializer to the registry.
func Register(id string, initializer func(cfg *config.Config) Agent) {
	if _, exists := agentInitializers[id]; exists {
		panic(fmt.Sprintf("agent already registered: %s", id))
	}
	agentInitializers[id] = initializer
}

// GetAgent initializes and returns an agent by its ID.
func GetAgent(id string, cfg *config.Config) (Agent, error) {
	initializer, ok := agentInitializers[id]
	if !ok {
		return nil, fmt.Errorf("agent not found: %s", id)
	}
	return initializer(cfg), nil
}

// ListAgents returns a list of available agent IDs.
func ListAgents() []string {
	var ids []string
	for id := range agentInitializers {
		ids = append(ids, id)
	}
	return ids
}
