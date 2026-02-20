package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/andybarilla/skeeter/internal/resolve"
	"github.com/andybarilla/skeeter/internal/store"
	"github.com/spf13/cobra"
)

var (
	listStatus   string
	listPriority string
	listAssignee string
	listTag      string
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List tasks",
	Aliases: []string{"ls"},
	RunE: func(cmd *cobra.Command, args []string) error {
		dir, err := resolve.Dir(dirFlag)
		if err != nil {
			return err
		}

		s, err := store.NewFilesystem(dir)
		if err != nil {
			return err
		}

		filter := store.Filter{
			Status:   listStatus,
			Priority: listPriority,
			Assignee: listAssignee,
			Tag:      listTag,
		}

		tasks, err := s.List(filter)
		if err != nil {
			return err
		}

		if len(tasks) == 0 {
			fmt.Println("No tasks found.")
			return nil
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tTITLE\tSTATUS\tPRIORITY\tASSIGNEE")
		for _, t := range tasks {
			title := t.Title
			if len(title) > 40 {
				title = title[:37] + "..."
			}
			assignee := t.Assignee
			if assignee == "" {
				assignee = "-"
			}
			priority := t.Priority
			if priority == "" {
				priority = "-"
			}
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", t.ID, title, t.Status, priority, assignee)
		}
		w.Flush()
		return nil
	},
}

func init() {
	listCmd.Flags().StringVar(&listStatus, "status", "", "filter by status")
	listCmd.Flags().StringVar(&listPriority, "priority", "", "filter by priority")
	listCmd.Flags().StringVar(&listAssignee, "assignee", "", "filter by assignee")
	listCmd.Flags().StringVar(&listTag, "tag", "", "filter by tag")
	rootCmd.AddCommand(listCmd)
}
