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
	"github.com/mauleyzaola/jira-etl/ports"
)

type changelogPage struct {
	StartAt    int               `json:"startAt"`
	MaxResults int               `json:"maxResults"`
	Total      int               `json:"total"`
	IsLast     bool              `json:"isLast"`
	Values     []domains.History `json:"values"`
}

var issueKeyPattern = regexp.MustCompile(`^[A-Z][A-Z0-9_]+-[0-9]+$`)
var projectKeyPattern = regexp.MustCompile(`^[A-Z][A-Z0-9_]+$`)

const filterDateFormat = "2006-01-02 15:04"

func JiraResponse(config envs.Config, req *http.Request) (*http.Response, error) {
	client := &http.Client{Timeout: 30 * time.Second}
	req.SetBasicAuth(config.JiraUserEmail, config.JiraApiToken)
	req.Header.Set("Accept", "application/json")
	return client.Do(req)
}

func IssueHistory(config envs.Config, key string) (*domains.JiraLog, error) {
	key = strings.ToUpper(key)
	if !issueKeyPattern.MatchString(key) {
		return nil, fmt.Errorf("invalid issue key: %q", key)
	}
	baseURL := config.JiraBaseUrl
	baseURL = strings.TrimRight(baseURL, "/")
	out := domains.JiraLog{Key: key, Histories: []domains.History{}}
	startAt := 0
	for {
		endpoint := fmt.Sprintf("%s/rest/api/3/issue/%s/changelog?startAt=%d&maxResults=200",
			baseURL, url.PathEscape(key), startAt)
		page, err := fetchPage(config, endpoint)
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

func IssueList(w io.Writer, formatter ports.Formatter, config envs.Config, project string, from, to time.Time) error {
	var (
		response = struct {
			Issues []domains.Issue `json:"issues"`
		}{}
	)
	project = strings.ToUpper(strings.TrimSpace(project))
	if !projectKeyPattern.MatchString(project) {
		return fmt.Errorf("invalid project key: %q", project)
	}
	baseURL := strings.TrimRight(config.JiraBaseUrl, "/")
	jql := fmt.Sprintf(`project = %s AND updated >= "%s" AND updated <= "%s" AND status != Cancelled ORDER BY updated ASC`,
		project,
		from.Format(filterDateFormat),
		to.Format(filterDateFormat),
	)
	endpoint := fmt.Sprintf("%s/rest/api/3/search/jql?jql=%s&maxResults=100&fields=key",
		baseURL, url.QueryEscape(jql))
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return err
	}
	resp, err := JiraResponse(config, req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("GET %s failed with HTTP %d: %s", endpoint, resp.StatusCode, body)
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return err
	}
	return formatter.Format(w, response.Issues)
}

func fetchPage(config envs.Config, endpoint string) (*changelogPage, error) {
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	resp, err := JiraResponse(config, req)
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
