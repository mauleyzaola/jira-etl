CREATE TABLE history
(
    id        BIGSERIAL PRIMARY KEY,
    issue_key TEXT        NOT NULL,
    author    TEXT        NOT NULL,
    created   TIMESTAMPTZ NOT NULL
);

CREATE INDEX idx_history_issue_key
    ON history (issue_key);

CREATE TABLE history_item
(
    id BIGSERIAL PRIMARY KEY,
    history_id BIGINT NOT NULL REFERENCES history(id),
    field TEXT NOT NULL,
    from_string TEXT NOT NULL,
    to_string TEXT NOT NULL
);