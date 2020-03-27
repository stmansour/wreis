package util

// AIRAuthenticateResponse is the reply structure from Accord Directory
type AIRAuthenticateResponse struct {
	Status   string       `json:"status"`
	UID      int64        `json:"uid"`
	Username string       `json:"username"` // user's first or preferred name
	Name     string       `json:"Name"`
	ImageURL string       `json:"ImageURL"`
	Message  string       `json:"message"`
	Token    string       `json:"Token"`
	Expire   JSONDateTime `json:"Expire"` // DATETIMEFMT in this format "2006-01-02T15:04:00Z"
}
