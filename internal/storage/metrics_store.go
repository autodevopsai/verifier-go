package storage

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Metric struct {
	AgentID    string    `json:"agent_id"`
	Timestamp  time.Time `json:"timestamp"`
	TokensUsed int       `json:"tokens_used"`
	Cost       float64   `json:"cost"`
	Result     string    `json:"result"`
	DurationMs int64     `json:"duration_ms"`
}

type MetricsStore struct {
	metricsDir string
}

func NewMetricsStore() *MetricsStore {
	return &MetricsStore{
		metricsDir: filepath.Join(".verifier", "metrics"),
	}
}

func (s *MetricsStore) Record(metric Metric) error {
	if err := os.MkdirAll(s.metricsDir, 0755); err != nil {
		return err
	}
	date := metric.Timestamp.Format("2006-01-02")
	filePath := filepath.Join(s.metricsDir, date+".json")

	var metrics []Metric
	data, err := os.ReadFile(filePath)
	if err == nil {
		_ = json.Unmarshal(data, &metrics)
	}

	metrics = append(metrics, metric)
	newData, err := json.MarshalIndent(metrics, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, newData, 0644)
}

func (s *MetricsStore) GetMetrics(period time.Duration) ([]Metric, error) {
	var results []Metric
	startTime := time.Now().Add(-period)

	files, err := os.ReadDir(s.metricsDir)
	if err != nil {
		return nil, nil // Return empty if dir doesn't exist
	}

	for _, file := range files {
		dateStr := strings.TrimSuffix(file.Name(), ".json")
		fileDate, err := time.Parse("2006-01-02", dateStr)
		if err != nil || fileDate.Before(startTime.Add(-24*time.Hour)) {
			continue // Parse error or file is too old
		}

		filePath := filepath.Join(s.metricsDir, file.Name())
		data, err := os.ReadFile(filePath)
		if err != nil {
			continue
		}

		var metrics []Metric
		if json.Unmarshal(data, &metrics) == nil {
			for _, m := range metrics {
				if m.Timestamp.After(startTime) {
					results = append(results, m)
				}
			}
		}
	}
	return results, nil
}
