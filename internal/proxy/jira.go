package proxy

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"
)

type JiraProxy struct{}

func NewJiraProxy() *JiraProxy {}

type changelogPage struct {
	StartAt    int               `json:"startAt"`
	MaxResults int               `json:"maxResults"`
	Total      int               `json:"total"`
	IsLast     bool              `json:"isLast"`
	Values     []json.RawMessage `json:"values"`
}

type output struct {
	Key       string            `json:"key"`
	Total     int               `json:"total"`
	Histories []json.RawMessage `json:"histories"`
}

var issueKeyPattern = regexp.MustCompile(`^[A-Z][A-Z0-9_]+-[0-9]+$`)

func run() error {
	if len(os.Args) != 2 {
		return fmt.Errorf("usage: %s ISSUE-KEY (e.g. DEV-6070)", os.Args[0])
	}
	key := strings.ToUpper(os.Args[1])
	if !issueKeyPattern.MatchString(key) {
		return fmt.Errorf("invalid issue key: %q", os.Args[1])
	}

	if err := loadDotEnv(".env"); err != nil {
		return err
	}
	baseURL, err := requireEnv("JIRA_BASE_URL")
	if err != nil {
		return err
	}
	email, err := requireEnv("JIRA_USER_EMAIL")
	if err != nil {
		return err
	}
	token, err := requireEnv("JIRA_API_TOKEN")
	if err != nil {
		return err
	}
	baseURL = strings.TrimRight(baseURL, "/")

	client := &http.Client{Timeout: 30 * time.Second}
	out := output{Key: key, Histories: []json.RawMessage{}}
	startAt := 0
	for {
		endpoint := fmt.Sprintf("%s/rest/api/3/issue/%s/changelog?startAt=%d&maxResults=100",
			baseURL, url.PathEscape(key), startAt)
		page, err := fetchPage(client, endpoint, email, token)
		if err != nil {
			return err
		}
		out.Total = page.Total
		out.Histories = append(out.Histories, page.Values...)
		startAt += len(page.Values)
		if page.IsLast || len(page.Values) == 0 {
			break
		}
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}

func fetchPage(client *http.Client, endpoint, email, token string) (*changelogPage, error) {
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(email, token)
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET %s failed with HTTP %d: %s", endpoint, resp.StatusCode, body)
	}

	var page changelogPage
	if err := json.Unmarshal(body, &page); err != nil {
		return nil, fmt.Errorf("decoding changelog: %w", err)
	}
	return &page, nil
}

// loadDotEnv sets KEY=VALUE pairs from path without overriding variables
// already present in the environment. A missing file is not an error.
func loadDotEnv(path string) error {
	f, err := os.Open(path)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		name, value, ok := strings.Cut(strings.TrimPrefix(line, "export "), "=")
		if !ok {
			continue
		}
		name = strings.TrimSpace(name)
		value = strings.Trim(strings.TrimSpace(value), `"'`)
		if _, exists := os.LookupEnv(name); !exists {
			os.Setenv(name, value)
		}
	}
	return scanner.Err()
}

func requireEnv(name string) (string, error) {
	value := os.Getenv(name)
	if value == "" {
		return "", fmt.Errorf("missing required environment variable: %s", name)
	}
	return value, nil
}
