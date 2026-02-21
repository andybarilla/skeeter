package main

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/andybarilla/skeeter/internal/store"
	"github.com/andybarilla/skeeter/internal/task"
	"github.com/spf13/cobra"
)

var (
	searchTitleOnly bool
	searchTag       string
)

var searchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Search tasks by title and body",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		s, err := openStore()
		if err != nil {
			return err
		}

		query := strings.ToLower(args[0])

		filter := store.Filter{}
		if searchTag != "" {
			filter.Tag = searchTag
		}

		allTasks, err := s.List(filter)
		if err != nil {
			return err
		}

		var results []task.Task
		for _, t := range allTasks {
			if searchTitleOnly {
				if strings.Contains(strings.ToLower(t.Title), query) {
					results = append(results, t)
				}
			} else {
				titleMatch := strings.Contains(strings.ToLower(t.Title), query)
				bodyMatch := strings.Contains(strings.ToLower(t.Body), query)
				if titleMatch || bodyMatch {
					results = append(results, t)
				}
			}
		}

		if isJSONOutput() {
			return outputTasksJSON(results)
		}
		if isYAMLOutput() {
			return outputTasksYAML(results)
		}

		if len(results) == 0 {
			fmt.Println("No matching tasks found.")
			return nil
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tTITLE\tSTATUS\tPRIORITY\tASSIGNEE")
		for _, t := range results {
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
	searchCmd.Flags().BoolVar(&searchTitleOnly, "title-only", false, "search only task titles")
	searchCmd.Flags().StringVar(&searchTag, "tag", "", "filter by tag in addition to query")
	rootCmd.AddCommand(searchCmd)
}
