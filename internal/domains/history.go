package domains

import (
	"encoding/json"
	"fmt"
	"time"
)

const jiraCreatedLayout = "2006-01-02T15:04:05.000-0700"

type History struct {
	ID              string           `json:"id"`
	Author          Author           `json:"author"`
	Created         time.Time        `json:"created"`
	Items           []HistoryItem    `json:"items"`
	HistoryMetadata *HistoryMetadata `json:"historyMetadata,omitempty"`
}

func (h *History) UnmarshalJSON(data []byte) error {
	type history History
	var raw struct {
		history
		Created string `json:"created"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	created, err := time.Parse(jiraCreatedLayout, raw.Created)
	if err != nil {
		return fmt.Errorf("parsing created: %w", err)
	}
	*h = History(raw.history)
	h.Created = created
	return nil
}
