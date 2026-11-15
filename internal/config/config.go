package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// DefaultConfig returns a production-ready default configuration.
func DefaultConfig() *Config {
	return &Config{
		Version: "1.0",
		Server: ServerConfig{
			Port:                6060,
			CORSAllowedOrigins:  []string{"*"},
			ReadTimeoutSeconds:  30,
			WriteTimeoutSeconds: 60,
			DossiersDir:         "./examples/dossiers",
		},
		Provider: ProviderConfig{
			Default:     "mock",
			APIKeyEnv:   "OPENAI_API_KEY",
			Model:       "gpt-4o",
			Temperature: 0.1,
			MaxTokens:   2048,
		},
		Protocol: ProtocolConfig{
			MaxRounds: 3,
			Thresholds: ThresholdConfig{
				Approve: 0.70,
				Reject:  0.35,
			},
			StrictCitations:     true,
			MinCitationCoverage: 0.80,
		},
		Ledger: LedgerConfig{
			StoragePath:       "./ledgers",
			HashAlgorithm:     "sha256",
			EnableMerkleRoots: true,
		},
	}
}

// LoadConfig reads a YAML configuration file from path, falling back to defaults for missing fields.
func LoadConfig(path string) (*Config, error) {
	cfg := DefaultConfig()

	if path == "" {
		return cfg, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return nil, fmt.Errorf("reading config file %s: %w", path, err)
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("parsing YAML config: %w", err)
	}

	// Resolve API key from environment if specified
	if cfg.Provider.APIKeyEnv != "" && cfg.Provider.APIKey == "" {
		cfg.Provider.APIKey = os.Getenv(cfg.Provider.APIKeyEnv)
	}

	return cfg, nil
}
