package resolve

import (
	"fmt"
	"os"
	"path/filepath"
)

func Dir(flagDir string) (string, error) {
	// 1. CLI flag
	if flagDir != "" {
		return filepath.Abs(flagDir)
	}

	// 2. Environment variable
	if envDir := os.Getenv("SKEETER_DIR"); envDir != "" {
		return filepath.Abs(envDir)
	}

	// 3. Walk up from cwd to find .skeeter/
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return findRoot(cwd)
}

func DirForInit(flagDir string) (string, error) {
	if flagDir != "" {
		return filepath.Abs(flagDir)
	}
	if envDir := os.Getenv("SKEETER_DIR"); envDir != "" {
		return filepath.Abs(envDir)
	}
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return filepath.Join(cwd, ".skeeter"), nil
}

func findRoot(from string) (string, error) {
	dir := from
	for {
		candidate := filepath.Join(dir, ".skeeter")
		if info, err := os.Stat(candidate); err == nil && info.IsDir() {
			return candidate, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("no .skeeter directory found (run 'skeeter init' to create one)")
		}
		dir = parent
	}
}
