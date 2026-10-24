CREATE TABLE jira_webhook
(
    ts                    BIGINT PRIMARY KEY,
    webhook_event         TEXT  NOT NULL,
    issue_event_type_name TEXT  NOT NULL,
    user_data             JSONB NOT NULL,
    issue_data            JSONB NOT NULL,
    changelog_data        JSONB NOT NULL
);