package main

import (
	"fmt"
	"strings"

	"github.com/andybarilla/skeeter/internal/resolve"
	"github.com/andybarilla/skeeter/internal/store"
	"github.com/spf13/cobra"
)

var showCmd = &cobra.Command{
	Use:   "show <id>",
	Short: "Show task details",
	Args:  cobra.ExactArgs(1),
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

		return nil
	},
}

func init() {
	rootCmd.AddCommand(showCmd)
}
