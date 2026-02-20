package main

import (
	"fmt"
	"strings"

	"github.com/andybarilla/skeeter/internal/resolve"
	"github.com/andybarilla/skeeter/internal/store"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status <id> <new-status>",
	Short: "Change task status",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		dir, err := resolve.Dir(dirFlag)
		if err != nil {
			return err
		}

		s, err := store.NewFilesystem(dir)
		if err != nil {
			return err
		}

		newStatus := args[1]
		if !s.Config.ValidStatus(newStatus) {
			return fmt.Errorf("invalid status %q (valid: %s)", newStatus, strings.Join(s.Config.Statuses, ", "))
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
