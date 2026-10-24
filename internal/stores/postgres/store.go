package postgres

import (
	"context"

	"github.com/mauleyzaola/jira-etl/internal/domains"
	"github.com/mauleyzaola/jira-etl/internal/ports"
)

type Store struct{}

func NewStore() *Store {
	return &Store{}
}

func (*Store) InsertJiraWebhook(ctx context.Context, db ports.DB, p domains.JiraWebhook) error {
	// assuming two timestamps won't collide, precision is milliseconds
	query := `
INSERT INTO jira_webhook(ts, webhook_event, issue_event_type_name, user_data, issue_data, changelog_data)
VALUES ($1, $2, $3, $4, $5, $6)
ON CONFLICT DO NOTHING;
`
	_, err := db.Exec(ctx, query, p.TS, p.WebhookEvent, p.IssueEventTypeName, p.User, p.Issue, p.ChangelogData)
	return err
}
