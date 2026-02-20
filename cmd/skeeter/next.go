package main

import (
	"fmt"
	"os"
	"sort"

	"github.com/andybarilla/skeeter/internal/resolve"
	"github.com/andybarilla/skeeter/internal/store"
	"github.com/andybarilla/skeeter/internal/task"
	"github.com/spf13/cobra"
)

var (
	nextAssign string
	nextQuiet  bool
)

var nextCmd = &cobra.Command{
	Use:   "next",
	Short: "Show the next available task for an agent to pick up",
	Long:  "Returns the highest-priority unassigned task in ready-for-development status. Designed for coding agents to discover work.",
	RunE: func(cmd *cobra.Command, args []string) error {
		dir, err := resolve.Dir(dirFlag)
		if err != nil {
			return err
		}

		s, err := store.NewFilesystem(dir)
		if err != nil {
			return err
		}

		tasks, err := s.List(store.Filter{Status: "ready-for-development"})
		if err != nil {
			return err
		}

		// Filter to unassigned only
		var available []task.Task
		for _, t := range tasks {
			if t.Assignee == "" {
				available = append(available, t)
			}
		}

		if len(available) == 0 {
			fmt.Fprintln(os.Stderr, "No tasks available")
			os.Exit(1)
		}

		// Sort by priority rank (lower = higher priority), then by ID
		sort.Slice(available, func(i, j int) bool {
			ri := s.Config.PriorityRank(available[i].Priority)
			rj := s.Config.PriorityRank(available[j].Priority)
			if ri != rj {
				return ri < rj
			}
			return available[i].ID < available[j].ID
		})

		picked := &available[0]

		if nextAssign != "" {
			picked.Assignee = nextAssign
			picked.Status = "in-progress"
			if err := s.Update(picked); err != nil {
				return err
			}
		}

		if nextQuiet {
			fmt.Println(picked.ID)
			return nil
		}

		printTask(picked)
		return nil
	},
}

func init() {
	nextCmd.Flags().StringVar(&nextAssign, "assign", "", "auto-assign the task and move to in-progress")
	nextCmd.Flags().BoolVarP(&nextQuiet, "quiet", "q", false, "output only the task ID")
	rootCmd.AddCommand(nextCmd)
}
