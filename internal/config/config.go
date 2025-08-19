package config

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

// Config corresponds to the structure of .verifier/config.yaml
type Config struct {
	Models    Models             `mapstructure:"models"`
	Providers Providers          `mapstructure:"providers"`
	Budgets   Budgets            `mapstructure:"budgets"`
	Thresholds Thresholds         `mapstructure:"thresholds"`
	Hooks     map[string][]string `mapstructure:"hooks"`
}

type Models struct {
	Primary  string `mapstructure:"primary"`
	Fallback string `mapstructure:"fallback"`
}

type Providers struct {
	OpenAI    ProviderAPIKey `mapstructure:"openai"`
	Anthropic ProviderAPIKey `mapstructure:"anthropic"`
}

type ProviderAPIKey struct {
	APIKey string `mapstructure:"api_key"`
}

type Budgets struct {
	DailyTokens     int `mapstructure:"daily_tokens"`
	PerCommitTokens int `mapstructure:"per_commit_tokens"`
	MonthlyCost     int `mapstructure:"monthly_cost"`
}

type Thresholds struct {
	DriftScore    int `mapstructure:"drift_score"`
	SecurityRisk  int `mapstructure:"security_risk"`
	CoverageDelta int `mapstructure:"coverage_delta"`
}

// Load reads configuration using Viper, respecting files, env vars, and .env
func Load() (*Config, error) {
	_ = godotenv.Load(filepath.Join(".verifier", ".env"))

	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".verifier")

	// For environment variables like VERIFIER_PROVIDERS_OPENAI_API_KEY
	v.SetEnvPrefix("VERIFIER")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	var cfg Config

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	// Manually bind from common env vars if not set by viper
	if cfg.Providers.OpenAI.APIKey == "" {
		cfg.Providers.OpenAI.APIKey = os.Getenv("OPENAI_API_KEY")
	}
	if cfg.Providers.Anthropic.APIKey == "" {
		cfg.Providers.Anthropic.APIKey = os.Getenv("ANTHROPIC_API_KEY")
	}

	return &cfg, nil
}

// Save writes the configuration to .verifier/config.yaml
func Save(cfg *Config) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	configPath := filepath.Join(".verifier", "config.yaml")
	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		return err
	}
	return os.WriteFile(configPath, data, 0644)
}
