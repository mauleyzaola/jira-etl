package proxy

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/mauleyzaola/jira-etl/domains"
	"github.com/mauleyzaola/jira-etl/envs"
)

type changelogPage struct {
	StartAt    int               `json:"startAt"`
	MaxResults int               `json:"maxResults"`
	Total      int               `json:"total"`
	IsLast     bool              `json:"isLast"`
	Values     []domains.History `json:"values"`
}

var issueKeyPattern = regexp.MustCompile(`^[A-Z][A-Z0-9_]+-[0-9]+$`)

func LoadJiraLog(config envs.Config, key string) (*domains.JiraLog, error) {
	key = strings.ToUpper(key)
	if !issueKeyPattern.MatchString(key) {
		return nil, fmt.Errorf("invalid issue key: %q", key)
	}
	baseURL := config.JiraBaseUrl
	email := config.JiraUserEmail
	token := config.JiraApiToken
	baseURL = strings.TrimRight(baseURL, "/")
	client := &http.Client{Timeout: 30 * time.Second}
	out := domains.JiraLog{Key: key, Histories: []domains.History{}}
	startAt := 0
	for {
		endpoint := fmt.Sprintf("%s/rest/api/3/issue/%s/changelog?startAt=%d&maxResults=200",
			baseURL, url.PathEscape(key), startAt)
		page, err := fetchPage(client, endpoint, email, token)
		if err != nil {
			return nil, err
		}
		out.Total = page.Total
		out.Histories = append(out.Histories, page.Values...)
		startAt += len(page.Values)
		if page.IsLast || len(page.Values) == 0 {
			break
		}
	}
	return &out, nil
}

func IssueList(config envs.Config, from, to time.Time) ([]byte, error) {
	baseURL := strings.TrimRight(config.JiraBaseUrl, "/")
	jql := fmt.Sprintf(`updated >= "%s" AND updated <= "%s" ORDER BY updated ASC`,
		from.Format("2006-01-02 15:04"),
		to.Format("2006-01-02 15:04"),
	)
	endpoint := fmt.Sprintf("%s/rest/api/3/search/jql?jql=%s&maxResults=100",
		baseURL, url.QueryEscape(jql))

	client := &http.Client{Timeout: 30 * time.Second}
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(config.JiraUserEmail, config.JiraApiToken)
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET %s failed with HTTP %d: %s", endpoint, resp.StatusCode, body)
	}
	return body, nil
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
	defer func() { _ = resp.Body.Close() }()

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
