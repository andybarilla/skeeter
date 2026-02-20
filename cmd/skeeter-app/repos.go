package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

type RepoEntry struct {
	Name   string `json:"name"`
	Path   string `json:"path"`
	Remote string `json:"remote"`
	Dir    string `json:"dir"`
}

type repoList struct {
	Repos []RepoEntry `json:"repos"`
}

type RepoStore struct {
	mu   sync.Mutex
	path string
}

func NewRepoStore() (*RepoStore, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}
	dir := filepath.Join(configDir, "skeeter-app")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}
	return &RepoStore{path: filepath.Join(dir, "repos.json")}, nil
}

func (rs *RepoStore) Load() ([]RepoEntry, error) {
	rs.mu.Lock()
	defer rs.mu.Unlock()
	return rs.loadLocked()
}

func (rs *RepoStore) loadLocked() ([]RepoEntry, error) {
	data, err := os.ReadFile(rs.path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var rl repoList
	if err := json.Unmarshal(data, &rl); err != nil {
		return nil, err
	}
	return rl.Repos, nil
}

func (rs *RepoStore) saveLocked(repos []RepoEntry) error {
	data, err := json.MarshalIndent(repoList{Repos: repos}, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(rs.path, data, 0644)
}

func (rs *RepoStore) Add(entry RepoEntry) error {
	rs.mu.Lock()
	defer rs.mu.Unlock()
	repos, err := rs.loadLocked()
	if err != nil {
		return err
	}
	for _, r := range repos {
		if r.Name == entry.Name {
			return fmt.Errorf("repo %q already exists", entry.Name)
		}
	}
	repos = append(repos, entry)
	return rs.saveLocked(repos)
}

func (rs *RepoStore) Remove(name string) error {
	rs.mu.Lock()
	defer rs.mu.Unlock()
	repos, err := rs.loadLocked()
	if err != nil {
		return err
	}
	for i, r := range repos {
		if r.Name == name {
			repos = append(repos[:i], repos[i+1:]...)
			return rs.saveLocked(repos)
		}
	}
	return fmt.Errorf("repo %q not found", name)
}

func (rs *RepoStore) MoveToFront(name string) error {
	rs.mu.Lock()
	defer rs.mu.Unlock()
	repos, err := rs.loadLocked()
	if err != nil {
		return err
	}
	for i, r := range repos {
		if r.Name == name {
			repos = append([]RepoEntry{r}, append(repos[:i], repos[i+1:]...)...)
			return rs.saveLocked(repos)
		}
	}
	return fmt.Errorf("repo %q not found", name)
}
