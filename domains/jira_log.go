package domains

type JiraLog struct {
	Key       string    `json:"key"`
	Total     int       `json:"total"`
	Histories []History `json:"histories"`
}
