package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/andybarilla/skeeter/internal/llm"
	"github.com/andybarilla/skeeter/internal/resolve"
	"github.com/spf13/cobra"
)

var (
	workMax    int
	workAssign string
	workDryRun bool
)

var workCmd = &cobra.Command{
	Use:   "work",
	Short: "Autonomous loop: pick tasks and invoke the LLM to implement them",
	Long: `Run an autonomous coding loop (the "Ralph Wiggum" technique).

Each iteration:
  1. Finds the highest-priority unassigned task in ready-for-development
  2. Claims it (assigns + moves to in-progress)
  3. Builds a prompt with task details
  4. Pipes the prompt to the configured LLM tool
  5. On success, marks the task done
  6. Repeats until no tasks remain or --max iterations reached

Configure the LLM tool:
  skeeter config set llm.tool claude`,
	RunE: func(cmd *cobra.Command, args []string) error {
		s, err := openStore()
		if err != nil {
			return err
		}

		cfg := s.GetConfig()

		tool, err := cfg.ResolveTool()
		if err != nil {
			return err
		}

		dir, err := resolve.Dir(dirFlag)
		if err != nil {
			return err
		}

		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer stop()

		iteration := 0
		for {
			if ctx.Err() != nil {
				fmt.Fprintln(os.Stderr, "\nInterrupted, stopping work loop.")
				return nil
			}

			if workMax > 0 && iteration >= workMax {
				fmt.Fprintf(os.Stderr, "Reached max iterations (%d), stopping.\n", workMax)
				return nil
			}

			picked, err := pickNextTask(s, cfg)
			if err != nil {
				return err
			}
			if picked == nil {
				fmt.Fprintln(os.Stderr, "No more tasks available, stopping.")
				return nil
			}

			fmt.Printf("\n=== Iteration %d: %s — %s ===\n\n", iteration+1, picked.ID, picked.Title)

			// Claim the task
			picked.Assignee = workAssign
			picked.Status = "in-progress"
			if err := s.Update(picked); err != nil {
				return err
			}

			systemPrompt, userContent := llm.BuildWorkPrompts(cfg, picked, dir)

			if workDryRun {
				fmt.Println("=== System Prompt ===")
				fmt.Println(systemPrompt)
				fmt.Println("\n=== User Content ===")
				fmt.Println(userContent)
				// Unclaim
				picked.Assignee = ""
				picked.Status = "ready-for-development"
				_ = s.Update(picked)
				return nil
			}

			// Execute work command
			if err := llm.RunCLIPassthrough(ctx, tool, systemPrompt, userContent, cfg.LLM.WorkArgs...); err != nil {
				fmt.Fprintf(os.Stderr, "\nWork command failed for %s: %v\n", picked.ID, err)
				// Revert task
				picked.Assignee = ""
				picked.Status = "ready-for-development"
				_ = s.Update(picked)
				return fmt.Errorf("work command failed for %s: %w", picked.ID, err)
			}

			// Re-read task from disk in case agent modified it
			fresh, err := s.Get(picked.ID)
			if err != nil {
				return fmt.Errorf("re-reading task %s: %w", picked.ID, err)
			}

			fresh.Status = "done"
			if err := s.Update(fresh); err != nil {
				return fmt.Errorf("marking %s done: %w", picked.ID, err)
			}

			fmt.Printf("\n=== %s marked done ===\n", picked.ID)
			iteration++
		}
	},
}

func init() {
	workCmd.Flags().IntVar(&workMax, "max", 0, "max iterations (0 = unlimited)")
	workCmd.Flags().StringVar(&workAssign, "assign", "ralph", "assignee name for claimed tasks")
	workCmd.Flags().BoolVar(&workDryRun, "dry-run", false, "print the prompt for the first task, then exit")
	rootCmd.AddCommand(workCmd)
}
