package main

import (
	"fmt"
	"strings"

	"github.com/andybarilla/skeeter/internal/resolve"
	"github.com/andybarilla/skeeter/internal/store"
	"github.com/spf13/cobra"
)

var assignCmd = &cobra.Command{
	Use:   "assign <id> <assignee>",
	Short: "Assign a task",
	Long:  "Assign a task to a person or agent. Use empty string to unassign.",
	Args:  cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		dir, err := resolve.Dir(dirFlag)
		if err != nil {
			return err
		}

		s, err := store.NewFilesystem(dir)
		if err != nil {
			return err
		}

		t, err := s.Get(strings.ToUpper(args[0]))
		if err != nil {
			return err
		}

		assignee := ""
		if len(args) > 1 {
			assignee = args[1]
		}

		t.Assignee = assignee

		if err := s.Update(t); err != nil {
			return err
		}

		if assignee == "" {
			fmt.Printf("%s: unassigned\n", t.ID)
		} else {
			fmt.Printf("%s: assigned to %s\n", t.ID, assignee)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(assignCmd)
}
