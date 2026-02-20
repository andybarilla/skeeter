package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/andybarilla/skeeter/internal/task"
	"github.com/spf13/cobra"
)

var (
	createPriority string
	createTags     string
	createAssignee string
	createStatus   string
)

var createCmd = &cobra.Command{
	Use:   "create <title>",
	Short: "Create a new task",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		s, err := openStore()
		if err != nil {
			return err
		}

		cfg := s.GetConfig()

		if createPriority != "" && !cfg.ValidPriority(createPriority) {
			return fmt.Errorf("invalid priority %q (valid: %s)", createPriority, strings.Join(cfg.Priorities, ", "))
		}

		if createStatus == "" {
			createStatus = cfg.Statuses[0]
		} else if !cfg.ValidStatus(createStatus) {
			return fmt.Errorf("invalid status %q (valid: %s)", createStatus, strings.Join(cfg.Statuses, ", "))
		}

		id, err := s.NextID()
		if err != nil {
			return err
		}

		now := time.Now().Format("2006-01-02")

		var tags task.FlowSlice
		if createTags != "" {
			for _, t := range strings.Split(createTags, ",") {
				tags = append(tags, strings.TrimSpace(t))
			}
		}

		t := &task.Task{
			ID:       id,
			Title:    args[0],
			Status:   createStatus,
			Priority: createPriority,
			Assignee: createAssignee,
			Tags:     tags,
			Created:  now,
			Updated:  now,
			Body:     "",
		}

		if err := s.Create(t); err != nil {
			return err
		}

		fmt.Printf("Created %s: %s\n", t.ID, t.Title)
		return nil
	},
}

func init() {
	createCmd.Flags().StringVarP(&createPriority, "priority", "p", "", "task priority")
	createCmd.Flags().StringVarP(&createTags, "tags", "t", "", "comma-separated tags")
	createCmd.Flags().StringVarP(&createAssignee, "assignee", "a", "", "task assignee")
	createCmd.Flags().StringVarP(&createStatus, "status", "s", "", "initial status (default: first configured status)")
	rootCmd.AddCommand(createCmd)
}
