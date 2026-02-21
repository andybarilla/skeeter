package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/andybarilla/skeeter/internal/llm"
	"github.com/spf13/cobra"
)

var enhanceCmd = &cobra.Command{
	Use:   "enhance <id>",
	Short: "Enhance a task description using AI",
	Long:  "Send a task to the configured LLM provider to generate a fleshed-out description with acceptance criteria, context, and technical considerations.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		s, err := openStore()
		if err != nil {
			return err
		}

		id := strings.ToUpper(args[0])
		t, err := s.Get(id)
		if err != nil {
			return err
		}

		template, err := s.LoadTemplate("default")
		if err != nil {
			template = ""
		}

		cfg := s.GetConfig()
		fmt.Printf("Enhancing %s: %s...\n", t.ID, t.Title)

		enhanced, err := llm.EnhanceTask(context.Background(), cfg, t, template)
		if err != nil {
			return fmt.Errorf("enhance failed: %w", err)
		}

		t.Body = enhanced
		if err := s.Update(t); err != nil {
			return fmt.Errorf("saving task: %w", err)
		}

		fmt.Println()
		printTask(t)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(enhanceCmd)
}
