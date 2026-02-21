package main

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/andybarilla/skeeter/internal/store"
	"github.com/andybarilla/skeeter/internal/task"
	"github.com/spf13/cobra"
)

var (
	listStatus      string
	listPriority    string
	listAssignee    string
	listTag         string
	listBlocked     bool
	listOverdue     bool
	listDueThisWeek bool
)

var listCmd = &cobra.Command{
	Use:     "list",
	Short:   "List tasks",
	Aliases: []string{"ls"},
	RunE: func(cmd *cobra.Command, args []string) error {
		s, err := openStore()
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

		if listBlocked {
			allTasks, _ := s.List(store.Filter{})
			var blocked []task.Task
			for _, t := range tasks {
				if store.IsBlocked(s, &t, allTasks) {
					blocked = append(blocked, t)
				}
			}
			tasks = blocked
		}

		today := time.Now().Format("2006-01-02")
		weekFromNow := time.Now().AddDate(0, 0, 7).Format("2006-01-02")

		if listOverdue {
			var overdue []task.Task
			for _, t := range tasks {
				if t.Due != "" && t.Due < today {
					overdue = append(overdue, t)
				}
			}
			tasks = overdue
		}

		if listDueThisWeek {
			var dueSoon []task.Task
			for _, t := range tasks {
				if t.Due != "" && t.Due >= today && t.Due <= weekFromNow {
					dueSoon = append(dueSoon, t)
				}
			}
			tasks = dueSoon
		}

		if isJSONOutput() {
			return outputTasksJSON(tasks)
		}
		if isYAMLOutput() {
			return outputTasksYAML(tasks)
		}

		if len(tasks) == 0 {
			fmt.Println("No tasks found.")
			return nil
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tTITLE\tSTATUS\tPRIORITY\tASSIGNEE\tDUE")
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
			due := t.Due
			if due == "" {
				due = "-"
			}
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n", t.ID, title, t.Status, priority, assignee, due)
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
	listCmd.Flags().BoolVar(&listBlocked, "blocked", false, "show only tasks with unmet dependencies")
	listCmd.Flags().BoolVar(&listOverdue, "overdue", false, "show only tasks past their due date")
	listCmd.Flags().BoolVar(&listDueThisWeek, "due-this-week", false, "show only tasks due in the next 7 days")
	rootCmd.AddCommand(listCmd)
}
