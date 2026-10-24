package domains

import "encoding/json"

// JiraWebhook is the incoming payload from the request.
type JiraWebhook struct {
	TS                 int64           `json:"timestamp"`
	WebhookEvent       string          `json:"webhookEvent"`
	IssueEventTypeName string          `json:"issue_event_type_name"`
	User               json.RawMessage `json:"user"`
	Issue              json.RawMessage `json:"issue"`
	ChangelogData      json.RawMessage `json:"changelog"`
}

// NewJiraWebhook deserializes a webhook body.
func NewJiraWebhook(data []byte) (*JiraWebhook, error) {
	var result JiraWebhook
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
