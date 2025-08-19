package agent

import (
	"encoding/json"
	"fmt"
	"time"
	
	"github.com/autodevops/verifier-go/internal/config"
	"github.com/autodevops/verifier-go/internal/provider"
)

type SecurityScanAgent struct {
	BaseAgent
	cfg *config.Config
}

func NewSecurityScanAgent(cfg *config.Config) Agent {
	return &SecurityScanAgent{
		BaseAgent: BaseAgent{
			id:          "security-scan",
			description: "Scans code for security vulnerabilities",
			model:       cfg.Models.Primary,
		},
		cfg: cfg,
	}
}

type SecurityAnalysis struct {
	RiskScore       int             `json:"risk_score"`
	Vulnerabilities []Vulnerability `json:"vulnerabilities"`
	Summary         string          `json:"summary"`
}

type Vulnerability struct {
	Type           string `json:"type"`
	Severity       string `json:"severity"`
	Description    string `json:"description"`
	Location       string `json:"location"`
	Recommendation string `json:"recommendation"`
}

func (a *SecurityScanAgent) Execute(ctx AgentContext) (*AgentResult, error) {
	if ctx.Diff == "" {
		res := a.CreateResult(AgentResult{Status: "skipped", Error: "No diff available"})
		return &res, nil
	}

	p, err := provider.ProviderFactory(a.Model(), a.cfg)
	if err != nil {
		return nil, err
	}

	prompt := fmt.Sprintf("Analyze the following code diff for security vulnerabilities.\n\n%s\n\nRespond JSON with { \"risk_score\": 0, \"vulnerabilities\": [{\"type\":\"\",\"severity\":\"critical|high|medium|low\",\"description\":\"\",\"location\":\"\",\"recommendation\":\"\"}], \"summary\":\"\" }", ctx.Diff)
	systemPrompt := "You are a security expert analyzing code for vulnerabilities. Be thorough but avoid false positives."

	response, err := p.Complete(prompt, systemPrompt, true)
	if err != nil {
		return nil, fmt.Errorf("security scan failed: %w", err)
	}
	
	var analysis SecurityAnalysis
	if err := json.Unmarshal([]byte(response), &analysis); err != nil {
		// If JSON fails, treat the whole response as a summary
		analysis.RiskScore = 3 // Default risk for unparsable response
		analysis.Summary = response
	}

	hasBlocking := false
	for _, v := range analysis.Vulnerabilities {
		if v.Severity == "critical" || v.Severity == "high" {
			hasBlocking = true
			break
		}
	}
	
	severity := "info"
	if hasBlocking {
		severity = "blocking"
	} else if analysis.RiskScore > 5 {
		severity = "warning"
	}

	// This is a rough estimation. A real implementation would get this from the provider's response.
	tokensUsed := len(prompt)/4 + len(response)/4
	
	res := a.CreateResult(AgentResult{
		Score:      analysis.RiskScore,
		Data:       analysis,
		Severity:   severity,
		TokensUsed: tokensUsed,
		Cost:       calculateCost(a.Model(), tokensUsed),
	})
	return &res, nil
}

func calculateCost(model string, tokens int) float64 {
	// Example costs per 1M tokens
	costPerMillion := 5.0 // GPT-4o
	if strings.Contains(model, "sonnet") {
		costPerMillion = 3.0 // Claude Sonnet 3.5
	}
	return (float64(tokens) / 1000000.0) * costPerMillion
}
