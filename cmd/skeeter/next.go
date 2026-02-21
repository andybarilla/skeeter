package main

import (
	"fmt"
	"os"

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
		s, err := openStore()
		if err != nil {
			return err
		}

		cfg := s.GetConfig()

		picked, err := pickNextTask(s, cfg)
		if err != nil {
			return err
		}

		if picked == nil {
			fmt.Fprintln(os.Stderr, "No tasks available")
			os.Exit(1)
		}

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
