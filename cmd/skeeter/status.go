package main

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status <id> <new-status>",
	Short: "Change task status",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		s, err := openStore()
		if err != nil {
			return err
		}

		cfg := s.GetConfig()
		newStatus := args[1]
		if !cfg.ValidStatus(newStatus) {
			return fmt.Errorf("invalid status %q (valid: %s)", newStatus, strings.Join(cfg.Statuses, ", "))
		}

		t, err := s.Get(strings.ToUpper(args[0]))
		if err != nil {
			return err
		}

		oldStatus := t.Status
		t.Status = newStatus

		if err := s.Update(t); err != nil {
			return err
		}

		fmt.Printf("%s: %s -> %s\n", t.ID, oldStatus, newStatus)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
