package agent

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type LintAgent struct {
	BaseAgent
}

func NewLintAgent() Agent {
	return &LintAgent{
		BaseAgent: BaseAgent{id: "lint", description: "Multi-language code linting", model: "none"},
	}
}

type LintIssue struct {
	File     string `json:"file"`
	Language string `json:"language"`
	Issues   string `json:"issues"`
}

func (a *LintAgent) Execute(ctx AgentContext) (*AgentResult, error) {
	if len(ctx.Files) == 0 {
		return &AgentResult{AgentID: a.ID(), Status: "skipped", Error: "No files to lint"}, nil
	}

	var allIssues []LintIssue
	totalIssues := 0

	for _, file := range ctx.Files {
		var lintOutput string
		var err error

		ext := filepath.Ext(file)
		lang := getLanguage(ext)

		switch ext {
		case ".go":
			lintOutput, err = runCommand("gofmt", "-l", file)
		case ".py":
			lintOutput, err = runCommand("ruff", "check", file)
		// Add cases for other languages like eslint for JS/TS
		default:
			continue
		}
		
		if err != nil {
			// Ruff exits with non-zero status if issues are found, which is not an execution error.
			if ext == ".py" && lintOutput != "" {
				// This is expected for ruff
			} else {
				fmt.Printf("linter for %s failed: %v, output: %s\n", file, err, lintOutput)
				continue
			}
		}

		if lintOutput != "" {
			allIssues = append(allIssues, LintIssue{File: file, Language: lang, Issues: lintOutput})
			totalIssues++
		}
	}

	artifactsDir := filepath.Join(".verifier", "artifacts")
	os.MkdirAll(artifactsDir, 0755)
	reportPath := filepath.Join(artifactsDir, "lint-report.json")

	if totalIssues > 0 {
		reportData, _ := json.MarshalIndent(allIssues, "", "  ")
		os.WriteFile(reportPath, reportData, 0644)
	}
	
	severity := "info"
	if totalIssues > 10 {
		severity = "warning"
	}

	result := a.CreateResult(AgentResult{
		Data: map[string]any{
			"total_issues":  totalIssues,
			"files_checked": len(ctx.Files),
			"issues":        allIssues,
		},
		Severity:  severity,
		Artifacts: []AgentArtifact{{Type: "report", Path: reportPath}},
	})
	return &result, nil
}

func runCommand(name string, args ...string) (string, error) {
	// Check if the command exists
	if _, err := exec.LookPath(name); err != nil {
		return fmt.Sprintf("%s not found in PATH", name), nil
	}
	cmd := exec.Command(name, args...)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

func getLanguage(ext string) string {
	switch ext {
	case ".go":
		return "Go"
	case ".py":
		return "Python"
	case ".js":
		return "JavaScript"
	case ".ts":
		return "TypeScript"
	default:
		return "Unknown"
	}
}
