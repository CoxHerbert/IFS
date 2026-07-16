package request

type AgentModelConfig struct {
	Label       string `json:"label"`
	Value       string `json:"value"`
	Description string `json:"description"`
	Default     bool   `json:"default"`
}

type AgentOllamaConfigRequest struct {
	BaseURL      string             `json:"baseUrl"`
	DefaultModel string             `json:"defaultModel"`
	Timeout      int                `json:"timeout"`
	Models       []AgentModelConfig `json:"models"`
}
