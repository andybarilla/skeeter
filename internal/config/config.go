package config

import (
	"os"
	"path/filepath"
	"slices"

	"gopkg.in/yaml.v3"
)

type ProjectConfig struct {
	Name   string `yaml:"name"`
	Prefix string `yaml:"prefix"`
}

type Config struct {
	Project    ProjectConfig `yaml:"project"`
	Statuses   []string      `yaml:"statuses"`
	Priorities []string      `yaml:"priorities"`
	AutoCommit bool          `yaml:"auto_commit"`
}

func Default() *Config {
	return &Config{
		Project: ProjectConfig{
			Name:   "",
			Prefix: "US",
		},
		Statuses:   []string{"backlog", "ready-for-development", "in-progress", "done"},
		Priorities: []string{"critical", "high", "medium", "low"},
		AutoCommit: false,
	}
}

func Load(dir string) (*Config, error) {
	path := filepath.Join(dir, "config.yaml")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	cfg := Default()
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

func (c *Config) Save(dir string) error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	path := filepath.Join(dir, "config.yaml")
	return os.WriteFile(path, data, 0644)
}

func (c *Config) ValidStatus(status string) bool {
	return slices.Contains(c.Statuses, status)
}

func (c *Config) ValidPriority(priority string) bool {
	return slices.Contains(c.Priorities, priority)
}
