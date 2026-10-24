package domains

type HistoryItem struct {
	Field            string  `json:"field"`
	FieldType        string  `json:"fieldtype"`
	FieldID          string  `json:"fieldId,omitempty"`
	From             *string `json:"from"`
	FromString       *string `json:"fromString"`
	To               *string `json:"to"`
	ToString         *string `json:"toString"`
	TmpFromAccountID *string `json:"tmpFromAccountId,omitempty"`
	TmpToAccountID   *string `json:"tmpToAccountId,omitempty"`
}
