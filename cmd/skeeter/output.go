package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/andybarilla/skeeter/internal/task"
	"gopkg.in/yaml.v3"
)

func outputTaskJSON(t *task.Task) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(t)
}

func outputTaskYAML(t *task.Task) error {
	enc := yaml.NewEncoder(os.Stdout)
	enc.SetIndent(2)
	return enc.Encode(t)
}

func outputTasksJSON(tasks []task.Task) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(tasks)
}

func outputTasksYAML(tasks []task.Task) error {
	enc := yaml.NewEncoder(os.Stdout)
	enc.SetIndent(2)
	return enc.Encode(tasks)
}

func outputNullJSON() error {
	fmt.Fprintln(os.Stdout, "null")
	return nil
}

func outputNullYAML() error {
	fmt.Fprintln(os.Stdout, "null")
	return nil
}
