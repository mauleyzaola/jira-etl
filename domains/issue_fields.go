package domains

import (
	"encoding/json"
	"fmt"
	"time"
)

type IssueFields struct {
	Summary   string    `json:"summary"`
	IssueType IssueType `json:"issuetype"`
	Status    Status    `json:"status"`
	Priority  *Priority `json:"priority"`
	Assignee  *Author   `json:"assignee"`
	Created   time.Time `json:"created"`
	Updated   time.Time `json:"updated"`
}

func (f *IssueFields) UnmarshalJSON(data []byte) error {
	type issueFields IssueFields
	var raw struct {
		issueFields
		Created string `json:"created"`
		Updated string `json:"updated"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	created, err := time.Parse(jiraCreatedLayout, raw.Created)
	if err != nil {
		return fmt.Errorf("parsing created: %w", err)
	}
	updated, err := time.Parse(jiraCreatedLayout, raw.Updated)
	if err != nil {
		return fmt.Errorf("parsing updated: %w", err)
	}
	*f = IssueFields(raw.issueFields)
	f.Created = created
	f.Updated = updated
	return nil
}
