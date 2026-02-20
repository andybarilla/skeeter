package id

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func Next(tasksDir, prefix string) (string, error) {
	entries, err := os.ReadDir(tasksDir)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Sprintf("%s-%03d", prefix, 1), nil
		}
		return "", err
	}

	maxNum := 0
	dashPrefix := prefix + "-"

	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := strings.TrimSuffix(e.Name(), filepath.Ext(e.Name()))
		if !strings.HasPrefix(name, dashPrefix) {
			continue
		}
		numStr := strings.TrimPrefix(name, dashPrefix)
		num, err := strconv.Atoi(numStr)
		if err != nil {
			continue
		}
		if num > maxNum {
			maxNum = num
		}
	}

	return fmt.Sprintf("%s-%03d", prefix, maxNum+1), nil
}
