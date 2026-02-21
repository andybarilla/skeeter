package main

import (
	"fmt"
	"strings"

	"github.com/andybarilla/skeeter/internal/resolve"
	"github.com/andybarilla/skeeter/internal/store"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "View or modify project configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		if remoteFlag != "" {
			return fmt.Errorf("config is not supported for remote repositories")
		}

		dir, err := resolve.Dir(dirFlag)
		if err != nil {
			return err
		}

		s, err := store.NewFilesystem(dir)
		if err != nil {
			return err
		}

		cfg := s.Config
		fmt.Printf("Project:     %s\n", cfg.Project.Name)
		fmt.Printf("Prefix:      %s\n", cfg.Project.Prefix)
		fmt.Printf("Statuses:    %s\n", strings.Join(cfg.Statuses, " -> "))
		fmt.Printf("Priorities:  %s\n", strings.Join(cfg.Priorities, ", "))
		fmt.Printf("Auto-commit: %v\n", cfg.AutoCommit)
		fmt.Printf("LLM command: %s\n", cfg.LLM.Command)
		return nil
	},
}

var configSetCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Set a configuration value",
	Long: `Set a configuration value. Available keys:
  name          Project name
  prefix        Task ID prefix (e.g., US, TASK, BUG)
  statuses      Comma-separated status list (ordered as workflow)
  priorities    Comma-separated priority list (highest first)
  auto_commit   Enable auto-commit (true/false)
  llm.command   LLM CLI command (e.g., "claude -p", "aichat -S", "llm")`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		if remoteFlag != "" {
			return fmt.Errorf("config is not supported for remote repositories")
		}

		dir, err := resolve.Dir(dirFlag)
		if err != nil {
			return err
		}

		s, err := store.NewFilesystem(dir)
		if err != nil {
			return err
		}

		key, value := args[0], args[1]

		switch key {
		case "name":
			s.Config.Project.Name = value
		case "prefix":
			s.Config.Project.Prefix = value
		case "statuses":
			var statuses []string
			for _, st := range strings.Split(value, ",") {
				statuses = append(statuses, strings.TrimSpace(st))
			}
			if len(statuses) < 2 {
				return fmt.Errorf("need at least 2 statuses")
			}
			s.Config.Statuses = statuses
		case "priorities":
			var priorities []string
			for _, p := range strings.Split(value, ",") {
				priorities = append(priorities, strings.TrimSpace(p))
			}
			if len(priorities) < 1 {
				return fmt.Errorf("need at least 1 priority")
			}
			s.Config.Priorities = priorities
		case "auto_commit":
			switch strings.ToLower(value) {
			case "true", "1", "yes":
				s.Config.AutoCommit = true
			case "false", "0", "no":
				s.Config.AutoCommit = false
			default:
				return fmt.Errorf("invalid value %q for auto_commit (use true/false)", value)
			}
		case "llm.command":
			s.Config.LLM.Command = value
		default:
			return fmt.Errorf("unknown config key %q (valid: name, prefix, statuses, priorities, auto_commit, llm.command)", key)
		}

		if err := s.Config.Save(dir); err != nil {
			return err
		}

		if err := s.RegenerateSkeeterMD(); err != nil {
			return err
		}

		fmt.Printf("Set %s = %s\n", key, value)
		return nil
	},
}

func init() {
	configCmd.AddCommand(configSetCmd)
	rootCmd.AddCommand(configCmd)
}
