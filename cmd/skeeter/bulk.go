package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var (
	bulkFromFile string
	bulkForce    bool
)

var bulkCmd = &cobra.Command{
	Use:   "bulk",
	Short: "Perform bulk operations on tasks",
}

var bulkStatusCmd = &cobra.Command{
	Use:   "status <status> [task-id]...",
	Short: "Change status of multiple tasks",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		s, err := openStore()
		if err != nil {
			return err
		}

		cfg := s.GetConfig()
		newStatus := args[0]
		if !cfg.ValidStatus(newStatus) {
			return fmt.Errorf("invalid status %q (valid: %s)", newStatus, strings.Join(cfg.Statuses, ", "))
		}

		ids, err := getTaskIDs(args[1:])
		if err != nil {
			return err
		}

		if len(ids) == 0 {
			fmt.Println("No task IDs provided.")
			return nil
		}

		if !bulkForce && len(ids) >= 5 {
			if !confirm(fmt.Sprintf("Change status to %q for %d tasks?", newStatus, len(ids))) {
				fmt.Println("Aborted.")
				return nil
			}
		}

		for _, id := range ids {
			t, err := s.Get(id)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Warning: %s: %v\n", id, err)
				continue
			}
			t.Status = newStatus
			if err := s.Update(t); err != nil {
				fmt.Fprintf(os.Stderr, "Error updating %s: %v\n", id, err)
				continue
			}
			fmt.Printf("%s: status -> %s\n", id, newStatus)
		}
		return nil
	},
}

var bulkAssignCmd = &cobra.Command{
	Use:   "assign <assignee> [task-id]...",
	Short: "Assign multiple tasks to someone",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		s, err := openStore()
		if err != nil {
			return err
		}

		assignee := args[0]
		ids, err := getTaskIDs(args[1:])
		if err != nil {
			return err
		}

		if len(ids) == 0 {
			fmt.Println("No task IDs provided.")
			return nil
		}

		if !bulkForce && len(ids) >= 5 {
			if !confirm(fmt.Sprintf("Assign %q to %d tasks?", assignee, len(ids))) {
				fmt.Println("Aborted.")
				return nil
			}
		}

		for _, id := range ids {
			t, err := s.Get(id)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Warning: %s: %v\n", id, err)
				continue
			}
			t.Assignee = assignee
			if err := s.Update(t); err != nil {
				fmt.Fprintf(os.Stderr, "Error updating %s: %v\n", id, err)
				continue
			}
			fmt.Printf("%s: assigned to %s\n", id, assignee)
		}
		return nil
	},
}

var bulkPriorityCmd = &cobra.Command{
	Use:   "priority <priority> [task-id]...",
	Short: "Change priority of multiple tasks",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		s, err := openStore()
		if err != nil {
			return err
		}

		cfg := s.GetConfig()
		newPriority := args[0]
		if !cfg.ValidPriority(newPriority) {
			return fmt.Errorf("invalid priority %q (valid: %s)", newPriority, strings.Join(cfg.Priorities, ", "))
		}

		ids, err := getTaskIDs(args[1:])
		if err != nil {
			return err
		}

		if len(ids) == 0 {
			fmt.Println("No task IDs provided.")
			return nil
		}

		if !bulkForce && len(ids) >= 5 {
			if !confirm(fmt.Sprintf("Set priority to %q for %d tasks?", newPriority, len(ids))) {
				fmt.Println("Aborted.")
				return nil
			}
		}

		for _, id := range ids {
			t, err := s.Get(id)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Warning: %s: %v\n", id, err)
				continue
			}
			t.Priority = newPriority
			if err := s.Update(t); err != nil {
				fmt.Fprintf(os.Stderr, "Error updating %s: %v\n", id, err)
				continue
			}
			fmt.Printf("%s: priority -> %s\n", id, newPriority)
		}
		return nil
	},
}

func getTaskIDs(args []string) ([]string, error) {
	ids := make([]string, 0)

	for _, arg := range args {
		ids = append(ids, strings.ToUpper(arg))
	}

	if bulkFromFile != "" {
		file, err := os.Open(bulkFromFile)
		if err != nil {
			return nil, fmt.Errorf("opening file: %w", err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line != "" && !strings.HasPrefix(line, "#") {
				ids = append(ids, strings.ToUpper(line))
			}
		}
		if err := scanner.Err(); err != nil {
			return nil, fmt.Errorf("reading file: %w", err)
		}
	}

	return ids, nil
}

func confirm(prompt string) bool {
	fmt.Printf("%s [y/N]: ", prompt)
	var response string
	fmt.Scanln(&response)
	return strings.ToLower(response) == "y" || strings.ToLower(response) == "yes"
}

func init() {
	bulkCmd.PersistentFlags().StringVar(&bulkFromFile, "from-file", "", "read task IDs from file (one per line)")
	bulkCmd.PersistentFlags().BoolVar(&bulkForce, "force", false, "skip confirmation prompt")

	bulkCmd.AddCommand(bulkStatusCmd)
	bulkCmd.AddCommand(bulkAssignCmd)
	bulkCmd.AddCommand(bulkPriorityCmd)

	rootCmd.AddCommand(bulkCmd)
}
