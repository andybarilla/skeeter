package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/andybarilla/skeeter/internal/resolve"
	"github.com/andybarilla/skeeter/internal/store"
	"github.com/spf13/cobra"
)

var (
	dirFlag    string
	remoteFlag string
	outputFlag string
)

var rootCmd = &cobra.Command{
	Use:   "skeeter",
	Short: "File-based project management for coding agents",
	Long:  "Skeeter is a file-based project management tool that stores tasks as markdown files in your git repository, designed for both humans and coding agents.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		if errors.Is(err, ErrNoTasksAvailable) {
			fmt.Fprintln(os.Stderr, "No tasks available")
		} else {
			fmt.Fprintln(os.Stderr, err)
		}
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&dirFlag, "dir", "", "path to skeeter directory (default: auto-detect .skeeter/)")
	rootCmd.PersistentFlags().StringVar(&remoteFlag, "remote", "", "use GitHub API backend (format: owner/repo)")
	rootCmd.PersistentFlags().StringVarP(&outputFlag, "output", "o", "table", "output format: table, json, yaml")
}

func openStore() (store.Store, error) {
	if remoteFlag != "" {
		return store.NewGitHub(remoteFlag, dirFlag)
	}
	dir, err := resolve.Dir(dirFlag)
	if err != nil {
		return nil, err
	}
	return store.NewFilesystem(dir)
}

func isJSONOutput() bool {
	return outputFlag == "json"
}

func isYAMLOutput() bool {
	return outputFlag == "yaml"
}

func isTableOutput() bool {
	return outputFlag == "table" || outputFlag == ""
}
