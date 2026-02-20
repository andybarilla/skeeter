package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/andybarilla/skeeter/internal/resolve"
	"github.com/andybarilla/skeeter/internal/store"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init [project-name]",
	Short: "Initialize skeeter in the current repository",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		dir, err := resolve.DirForInit(dirFlag)
		if err != nil {
			return err
		}

		if _, err := os.Stat(filepath.Join(dir, "config.yaml")); err == nil {
			return fmt.Errorf("skeeter already initialized at %s", dir)
		}

		projectName := ""
		if len(args) > 0 {
			projectName = args[0]
		} else {
			cwd, _ := os.Getwd()
			projectName = filepath.Base(cwd)
		}

		s := &store.FilesystemStore{Dir: dir}
		if err := s.Init(projectName); err != nil {
			return err
		}

		fmt.Printf("Initialized skeeter at %s\n", dir)
		fmt.Printf("  Project: %s\n", projectName)
		fmt.Printf("  Prefix:  %s\n", s.Config.Project.Prefix)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
