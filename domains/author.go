package domains

type Author struct {
	Self         string     `json:"self"`
	AccountID    string     `json:"accountId"`
	EmailAddress string     `json:"emailAddress,omitempty"`
	AvatarUrls   AvatarUrls `json:"avatarUrls"`
	DisplayName  string     `json:"displayName"`
	Active       bool       `json:"active"`
	TimeZone     string     `json:"timeZone"`
	AccountType  string     `json:"accountType"`
}
