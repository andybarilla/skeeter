package main

import (
	"fmt"
	"strings"

	"github.com/andybarilla/skeeter/internal/store"
	"github.com/andybarilla/skeeter/internal/task"
	"github.com/spf13/cobra"
)

var showCmd = &cobra.Command{
	Use:   "show <id>",
	Short: "Show task details",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		s, err := openStore()
		if err != nil {
			return err
		}

		t, err := s.Get(strings.ToUpper(args[0]))
		if err != nil {
			return err
		}

		if isJSONOutput() {
			return outputTaskJSON(t)
		}
		if isYAMLOutput() {
			return outputTaskYAML(t)
		}

		allTasks, _ := s.List(store.Filter{})
		printTaskWithDeps(t, s, allTasks)
		return nil
	},
}

func printTaskWithDeps(t *task.Task, s store.Store, allTasks []task.Task) {
	printTask(t)

	depStatus := store.GetDependencyStatus(s, t, allTasks)

	if len(depStatus.DependsOn) > 0 {
		fmt.Println()
		if len(depStatus.BlockedBy) > 0 {
			fmt.Printf("Blocked by: %s (incomplete)\n", strings.Join(depStatus.BlockedBy, ", "))
		} else {
			fmt.Printf("Depends on: %s (all complete)\n", strings.Join(depStatus.DependsOn, ", "))
		}
	}
	if len(depStatus.Blocking) > 0 {
		fmt.Printf("Blocking:   %s\n", strings.Join(depStatus.Blocking, ", "))
	}
}

func printTask(t *task.Task) {
	assignee := t.Assignee
	if assignee == "" {
		assignee = "-"
	}
	priority := t.Priority
	if priority == "" {
		priority = "-"
	}

	fmt.Printf("%s: %s\n", t.ID, t.Title)
	fmt.Printf("Status: %s | Priority: %s | Assignee: %s\n", t.Status, priority, assignee)
	if t.Due != "" {
		fmt.Printf("Due: %s\n", t.Due)
	}
	if len(t.Tags) > 0 {
		fmt.Printf("Tags: %s\n", strings.Join(t.Tags, ", "))
	}
	if len(t.Links) > 0 {
		fmt.Printf("Links: %s\n", strings.Join(t.Links, ", "))
	}
	fmt.Printf("Created: %s | Updated: %s\n", t.Created, t.Updated)

	if t.Body != "" {
		fmt.Println()
		fmt.Print(t.Body)
	}
}

func init() {
	rootCmd.AddCommand(showCmd)
}
