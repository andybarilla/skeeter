package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/andybarilla/skeeter/internal/resolve"
	"github.com/andybarilla/skeeter/internal/store"
	"github.com/spf13/cobra"
)

var editCmd = &cobra.Command{
	Use:   "edit <id>",
	Short: "Open a task in your editor",
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

		taskID := strings.ToUpper(args[0])

		// Verify it exists
		if _, err := s.Get(taskID); err != nil {
			return err
		}

		editor := os.Getenv("EDITOR")
		if editor == "" {
			editor = "vi"
		}

		taskPath := filepath.Join(dir, "tasks", taskID+".md")
		c := exec.Command(editor, taskPath)
		c.Stdin = os.Stdin
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr

		if err := c.Run(); err != nil {
			return fmt.Errorf("editor exited with error: %w", err)
		}

		// Re-parse to update the timestamp
		t, err := s.Get(taskID)
		if err != nil {
			return fmt.Errorf("warning: file may have invalid format: %w", err)
		}

		if err := s.Update(t); err != nil {
			return err
		}

		fmt.Printf("Updated %s\n", taskID)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(editCmd)
}
