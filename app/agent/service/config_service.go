package service

import (
	"baize/app/agent/dao"
	"baize/app/agent/request"
	"baize/app/setting"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const (
	agentConfigBaseURL      = "ollama.base_url"
	agentConfigDefaultModel = "ollama.default_model"
	agentConfigTimeout      = "ollama.timeout"
	agentConfigModels       = "ollama.models"
)

type ConfigService struct {
	Dao *dao.ConfigDao
}

var configService = &ConfigService{Dao: dao.GetConfigDao()}

func GetConfigService() *ConfigService {
	return configService
}

func (s *ConfigService) GetOllamaConfig() *request.AgentOllamaConfigRequest {
	cfg := yamlOllamaConfig()
	values := s.safeSelectAll()
	if value := strings.TrimSpace(values[agentConfigBaseURL]); value != "" {
		cfg.BaseURL = value
	}
	if value := strings.TrimSpace(values[agentConfigDefaultModel]); value != "" {
		cfg.DefaultModel = value
	}
	if value := strings.TrimSpace(values[agentConfigTimeout]); value != "" {
		var timeout int
		if err := json.Unmarshal([]byte(value), &timeout); err == nil && timeout > 0 {
			cfg.Timeout = timeout
		}
	}
	if value := strings.TrimSpace(values[agentConfigModels]); value != "" {
		var models []request.AgentModelConfig
		if err := json.Unmarshal([]byte(value), &models); err == nil && len(models) > 0 {
			cfg.Models = models
		}
	}
	normalizeOllamaConfig(cfg)
	return cfg
}

func (s *ConfigService) SaveOllamaConfig(cfg *request.AgentOllamaConfigRequest) *request.AgentOllamaConfigRequest {
	normalizeOllamaConfig(cfg)
	timeoutJSON, _ := json.Marshal(cfg.Timeout)
	modelsJSON, _ := json.Marshal(cfg.Models)
	s.Dao.Upsert(map[string]string{
		agentConfigBaseURL:      cfg.BaseURL,
		agentConfigDefaultModel: cfg.DefaultModel,
		agentConfigTimeout:      string(timeoutJSON),
		agentConfigModels:       string(modelsJSON),
	})
	return cfg
}

func (s *ConfigService) TestOllamaConfig(cfg *request.AgentOllamaConfigRequest) error {
	normalizeOllamaConfig(cfg)
	client := &http.Client{Timeout: time.Duration(cfg.Timeout) * time.Second}
	resp, err := client.Get(cfg.BaseURL + "/api/tags")
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("ollama returned status %d", resp.StatusCode)
	}
	return nil
}

func (s *ConfigService) safeSelectAll() map[string]string {
	defer func() {
		_ = recover()
	}()
	return s.Dao.SelectAll()
}

func yamlOllamaConfig() *request.AgentOllamaConfigRequest {
	cfg := &request.AgentOllamaConfigRequest{
		BaseURL:      "http://localhost:11434",
		DefaultModel: fallbackModelName,
		Timeout:      90,
		Models:       []request.AgentModelConfig{},
	}
	if setting.Conf.AgentConfig == nil || setting.Conf.AgentConfig.OllamaConfig == nil {
		return cfg
	}
	ollama := setting.Conf.AgentConfig.OllamaConfig
	if strings.TrimSpace(ollama.BaseURL) != "" {
		cfg.BaseURL = strings.TrimSpace(ollama.BaseURL)
	}
	if strings.TrimSpace(ollama.DefaultModel) != "" {
		cfg.DefaultModel = strings.TrimSpace(ollama.DefaultModel)
	}
	if ollama.Timeout > 0 {
		cfg.Timeout = ollama.Timeout
	}
	for _, item := range ollama.Models {
		cfg.Models = append(cfg.Models, request.AgentModelConfig{
			Label:       item.Label,
			Value:       item.Value,
			Description: item.Description,
			Default:     item.Default,
		})
	}
	return cfg
}

func normalizeOllamaConfig(cfg *request.AgentOllamaConfigRequest) {
	cfg.BaseURL = strings.TrimRight(strings.TrimSpace(cfg.BaseURL), "/")
	if cfg.BaseURL == "" {
		cfg.BaseURL = "http://localhost:11434"
	}
	if cfg.Timeout <= 0 {
		cfg.Timeout = 90
	}
	models := make([]request.AgentModelConfig, 0, len(cfg.Models))
	hasDefault := false
	for _, item := range cfg.Models {
		item.Label = strings.TrimSpace(item.Label)
		item.Value = strings.TrimSpace(item.Value)
		item.Description = strings.TrimSpace(item.Description)
		if item.Value == "" {
			continue
		}
		if item.Label == "" {
			item.Label = item.Value
		}
		if item.Default {
			hasDefault = true
			cfg.DefaultModel = item.Value
		}
		models = append(models, item)
	}
	if len(models) == 0 {
		models = []request.AgentModelConfig{{Label: "Qwen 2.5 7B", Value: fallbackModelName, Description: "默认模型", Default: true}}
		hasDefault = true
		cfg.DefaultModel = fallbackModelName
	}
	if strings.TrimSpace(cfg.DefaultModel) == "" {
		cfg.DefaultModel = models[0].Value
	}
	if !hasDefault {
		for i := range models {
			if models[i].Value == cfg.DefaultModel {
				models[i].Default = true
				hasDefault = true
				break
			}
		}
	}
	if !hasDefault {
		models[0].Default = true
		cfg.DefaultModel = models[0].Value
	}
	cfg.Models = models
}
