package config

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"

	"gopkg.in/yaml.v3"
)

type ProjectConfig struct {
	Name   string `yaml:"name" json:"name"`
	Prefix string `yaml:"prefix" json:"prefix"`
}

// LLMToolDef describes how to invoke an LLM CLI tool.
type LLMToolDef struct {
	Command          string `yaml:"command" json:"command"`
	PrintFlag        string `yaml:"print_flag" json:"print_flag"`
	SystemPromptFlag string `yaml:"system_prompt_flag" json:"system_prompt_flag"`
}

var builtinTools = map[string]LLMToolDef{
	"claude": {
		Command:          "claude",
		PrintFlag:        "-p",
		SystemPromptFlag: "--system-prompt",
	},
}

type LLMConfig struct {
	Tool     string                `yaml:"tool,omitempty" json:"tool"`
	Tools    map[string]LLMToolDef `yaml:"tools,omitempty" json:"tools,omitempty"`
	WorkArgs []string              `yaml:"work_args,omitempty" json:"work_args,omitempty"`
}

type Config struct {
	Project    ProjectConfig `yaml:"project" json:"project"`
	Statuses   []string      `yaml:"statuses" json:"statuses"`
	Priorities []string      `yaml:"priorities" json:"priorities"`
	AutoCommit bool          `yaml:"auto_commit" json:"auto_commit"`
	LLM        LLMConfig     `yaml:"llm,omitempty" json:"llm"`
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
		LLM:        LLMConfig{Tool: "claude"},
	}
}

// ResolveTool returns the LLMToolDef for the configured tool name.
// User-defined tools (llm.tools) take precedence over builtins.
func (c *Config) ResolveTool() (*LLMToolDef, error) {
	name := c.LLM.Tool
	if name == "" {
		return nil, fmt.Errorf("no LLM tool configured (run: skeeter config set llm.tool claude)")
	}
	if c.LLM.Tools != nil {
		if t, ok := c.LLM.Tools[name]; ok {
			return &t, nil
		}
	}
	if t, ok := builtinTools[name]; ok {
		return &t, nil
	}
	return nil, fmt.Errorf("unknown LLM tool %q (builtin: claude)", name)
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

// PriorityRank returns the index of a priority in the configured list.
// Lower index = higher priority. Returns len(Priorities) for unknown values.
func (c *Config) PriorityRank(priority string) int {
	idx := slices.Index(c.Priorities, priority)
	if idx < 0 {
		return len(c.Priorities)
	}
	return idx
}
