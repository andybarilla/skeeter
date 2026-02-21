package store

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/andybarilla/skeeter/internal/config"
	"github.com/andybarilla/skeeter/internal/id"
	"github.com/andybarilla/skeeter/internal/task"
	"gopkg.in/yaml.v3"
)

type GitHubStore struct {
	owner   string
	repo    string
	dir     string
	token   string
	client  *http.Client
	cfg     *config.Config
	baseURL string
}

type ghContentsResponse struct {
	Name     string `json:"name"`
	Path     string `json:"path"`
	SHA      string `json:"sha"`
	Content  string `json:"content"`
	Encoding string `json:"encoding"`
	Type     string `json:"type"`
}

func NewGitHub(remote, dir string) (*GitHubStore, error) {
	owner, repo, found := strings.Cut(remote, "/")
	if !found {
		return nil, fmt.Errorf("invalid remote format %q (expected owner/repo)", remote)
	}

	token, err := resolveToken()
	if err != nil {
		return nil, err
	}

	if dir == "" {
		dir = ".skeeter"
	}

	s := &GitHubStore{
		owner:  owner,
		repo:   repo,
		dir:    dir,
		token:  token,
		client: &http.Client{Timeout: 30 * time.Second},
	}

	cfg, err := s.loadConfig()
	if err != nil {
		return nil, fmt.Errorf("loading remote config: %w", err)
	}
	s.cfg = cfg

	return s, nil
}

func resolveToken() (string, error) {
	out, err := exec.Command("gh", "auth", "token").Output()
	if err == nil {
		token := strings.TrimSpace(string(out))
		if token != "" {
			return token, nil
		}
	}

	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		return token, nil
	}

	return "", fmt.Errorf("no GitHub token found (install gh CLI or set GITHUB_TOKEN)")
}

func (s *GitHubStore) contentsURL(path string) string {
	base := s.baseURL
	if base == "" {
		base = "https://api.github.com"
	}
	return fmt.Sprintf("%s/repos/%s/%s/contents/%s", base, s.owner, s.repo, path)
}

func (s *GitHubStore) tasksPath() string {
	return s.dir + "/tasks"
}

func (s *GitHubStore) taskFilePath(taskID string) string {
	return s.tasksPath() + "/" + taskID + ".md"
}

func (s *GitHubStore) doRequest(method, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+s.token)
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	return s.client.Do(req)
}

func (s *GitHubStore) getFileContent(path string) (content []byte, sha string, err error) {
	resp, err := s.doRequest("GET", s.contentsURL(path), nil)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return nil, "", fmt.Errorf("file not found: %s", path)
	}
	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, "", fmt.Errorf("GitHub API error %d: %s", resp.StatusCode, body)
	}

	var cr ghContentsResponse
	if err := json.NewDecoder(resp.Body).Decode(&cr); err != nil {
		return nil, "", fmt.Errorf("decoding response: %w", err)
	}

	decoded, err := base64.StdEncoding.DecodeString(strings.ReplaceAll(cr.Content, "\n", ""))
	if err != nil {
		return nil, "", fmt.Errorf("decoding base64 content: %w", err)
	}

	return decoded, cr.SHA, nil
}

func (s *GitHubStore) putFile(path string, content []byte, sha, message string) error {
	payload := map[string]string{
		"message": "skeeter: " + message,
		"content": base64.StdEncoding.EncodeToString(content),
	}
	if sha != "" {
		payload["sha"] = sha
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := s.doRequest("PUT", s.contentsURL(path), strings.NewReader(string(data)))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("GitHub API error %d: %s", resp.StatusCode, body)
	}

	return nil
}

func (s *GitHubStore) listDir(path string) ([]ghContentsResponse, error) {
	resp, err := s.doRequest("GET", s.contentsURL(path), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return nil, nil
	}
	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("GitHub API error %d: %s", resp.StatusCode, body)
	}

	var entries []ghContentsResponse
	if err := json.NewDecoder(resp.Body).Decode(&entries); err != nil {
		return nil, fmt.Errorf("decoding directory listing: %w", err)
	}

	return entries, nil
}

func (s *GitHubStore) loadConfig() (*config.Config, error) {
	data, _, err := s.getFileContent(s.dir + "/config.yaml")
	if err != nil {
		return nil, err
	}

	cfg := config.Default()
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

func (s *GitHubStore) GetConfig() *config.Config {
	return s.cfg
}

func (s *GitHubStore) Init(projectName string) error {
	return fmt.Errorf("init is not supported for remote repositories — initialize locally and push")
}

func (s *GitHubStore) List(filter Filter) ([]task.Task, error) {
	entries, err := s.listDir(s.tasksPath())
	if err != nil {
		return nil, err
	}

	var tasks []task.Task
	for _, e := range entries {
		if e.Type != "file" || filepath.Ext(e.Name) != ".md" {
			continue
		}

		data, _, err := s.getFileContent(e.Path)
		if err != nil {
			continue
		}

		t, err := task.Parse(string(data))
		if err != nil {
			continue
		}

		if !matchesFilter(t, filter) {
			continue
		}

		tasks = append(tasks, *t)
	}

	return tasks, nil
}

func (s *GitHubStore) Get(taskID string) (*task.Task, error) {
	data, _, err := s.getFileContent(s.taskFilePath(taskID))
	if err != nil {
		return nil, fmt.Errorf("task %s not found", taskID)
	}
	return task.Parse(string(data))
}

func (s *GitHubStore) Create(t *task.Task) error {
	content, err := task.Marshal(t)
	if err != nil {
		return err
	}

	return s.putFile(
		s.taskFilePath(t.ID),
		[]byte(content),
		"",
		fmt.Sprintf("create %s: %s", t.ID, t.Title),
	)
}

func (s *GitHubStore) Update(t *task.Task) error {
	t.Updated = time.Now().Format("2006-01-02")

	// Fetch current SHA for conflict detection
	_, sha, err := s.getFileContent(s.taskFilePath(t.ID))
	if err != nil {
		return fmt.Errorf("fetching current version of %s: %w", t.ID, err)
	}

	content, err := task.Marshal(t)
	if err != nil {
		return err
	}

	return s.putFile(
		s.taskFilePath(t.ID),
		[]byte(content),
		sha,
		fmt.Sprintf("update %s: %s", t.ID, t.Title),
	)
}

func (s *GitHubStore) NextID() (string, error) {
	entries, err := s.listDir(s.tasksPath())
	if err != nil {
		return "", err
	}

	// Build a temp list of filenames for id.Next-style logic
	var names []string
	for _, e := range entries {
		names = append(names, e.Name)
	}

	return id.NextFromNames(names, s.cfg.Project.Prefix)
}

func (s *GitHubStore) LoadTemplate(name string) (string, error) {
	path := s.dir + "/templates/" + name + ".md"
	data, _, err := s.getFileContent(path)
	if err != nil {
		return "", fmt.Errorf("template %q not found in remote repository", name)
	}
	return string(data), nil
}
