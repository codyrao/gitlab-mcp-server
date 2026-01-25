package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server  ServerConfig  `yaml:"server"`
	GitLab  GitLabConfig  `yaml:"gitlab"`
	Logging LoggingConfig `yaml:"logging"`
}

type ServerConfig struct {
	Host      string `yaml:"host"`
	Port      int    `yaml:"port"`
	Transport string `yaml:"transport"`
}

type GitLabConfig struct {
	Host  string `yaml:"host"`
	Token string `yaml:"token"`
}

type LoggingConfig struct {
	Level string `yaml:"level"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	if token := os.Getenv("GITLAB_TOKEN"); token != "" {
		config.GitLab.Token = token
	}
	if host := os.Getenv("GITLAB_HOST"); host != "" {
		config.GitLab.Host = host
	}
	if transport := os.Getenv("GITLAB_MCP_TRANSPORT"); transport != "" {
		config.Server.Transport = transport
	}
	if port := os.Getenv("GITLAB_MCP_PORT"); port != "" {
		var p int
		if _, err := fmt.Sscanf(port, "%d", &p); err == nil {
			config.Server.Port = p
		}
	}

	return &config, nil
}

func (c *Config) Validate() error {
	if c.GitLab.Host == "" {
		return fmt.Errorf("gitlab host is required")
	}
	if c.GitLab.Token == "" {
		return fmt.Errorf("gitlab token is required")
	}
	if c.Server.Transport == "" {
		c.Server.Transport = "stdio"
	}
	if c.Server.Port == 0 {
		c.Server.Port = 8080
	}
	return nil
}
