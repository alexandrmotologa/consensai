package config

// Config defines the complete configuration for the consensai engine.
type Config struct {
	Version  string         `yaml:"version" json:"version"`
	Server   ServerConfig   `yaml:"server" json:"server"`
	Provider ProviderConfig `yaml:"provider" json:"provider"`
	Protocol ProtocolConfig `yaml:"protocol" json:"protocol"`
	Ledger   LedgerConfig   `yaml:"ledger" json:"ledger"`
}

// ServerConfig defines the HTTP and API settings.
type ServerConfig struct {
	Port                int      `yaml:"port" json:"port"`
	CORSAllowedOrigins  []string `yaml:"cors_allowed_origins" json:"cors_allowed_origins"`
	ReadTimeoutSeconds  int      `yaml:"read_timeout_seconds" json:"read_timeout_seconds"`
	WriteTimeoutSeconds int      `yaml:"write_timeout_seconds" json:"write_timeout_seconds"`
	DossiersDir         string   `yaml:"dossiers_dir" json:"dossiers_dir"`
}

// ProviderConfig specifies the LLM provider configuration.
type ProviderConfig struct {
	Default    string  `yaml:"default" json:"default"` // mock, openai, anthropic, gemini, ollama
	APIKeyEnv  string  `yaml:"api_key_env" json:"api_key_env"`
	APIKey     string  `yaml:"api_key,omitempty" json:"api_key,omitempty"`
	BaseURL    string  `yaml:"base_url,omitempty" json:"base_url,omitempty"`
	Model      string  `yaml:"model" json:"model"`
	Temperature float64 `yaml:"temperature" json:"temperature"`
	MaxTokens  int     `yaml:"max_tokens" json:"max_tokens"`
}

// ProtocolConfig holds parameters for the 3-round dialectic arbitration.
type ProtocolConfig struct {
	MaxRounds            int               `yaml:"max_rounds" json:"max_rounds"`
	Thresholds           ThresholdConfig   `yaml:"thresholds" json:"thresholds"`
	StrictCitations      bool              `yaml:"strict_citations" json:"strict_citations"`
	MinCitationCoverage  float64           `yaml:"min_citation_coverage" json:"min_citation_coverage"`
}

// ThresholdConfig holds the decision boundaries for Epistemic Convergence.
type ThresholdConfig struct {
	Approve float64 `yaml:"approve" json:"approve"` // default 0.70
	Reject  float64 `yaml:"reject" json:"reject"`   // default 0.35
}

// LedgerConfig holds the audit trail parameters.
type LedgerConfig struct {
	StoragePath       string `yaml:"storage_path" json:"storage_path"`
	HashAlgorithm     string `yaml:"hash_algorithm" json:"hash_algorithm"`
	EnableMerkleRoots bool   `yaml:"enable_merkle_roots" json:"enable_merkle_roots"`
}
