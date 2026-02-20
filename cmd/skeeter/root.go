package main

import (
	"fmt"
	"os"

	"github.com/andybarilla/skeeter/internal/resolve"
	"github.com/andybarilla/skeeter/internal/store"
	"github.com/spf13/cobra"
)

var (
	dirFlag    string
	remoteFlag string
)

var rootCmd = &cobra.Command{
	Use:   "skeeter",
	Short: "File-based project management for coding agents",
	Long:  "Skeeter is a file-based project management tool that stores tasks as markdown files in your git repository, designed for both humans and coding agents.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&dirFlag, "dir", "", "path to skeeter directory (default: auto-detect .skeeter/)")
	rootCmd.PersistentFlags().StringVar(&remoteFlag, "remote", "", "use GitHub API backend (format: owner/repo)")
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
