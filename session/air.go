package session

import "mojo/util"

// AIRAuthenticateResponse is the reply structure from Accord Directory
type AIRAuthenticateResponse struct {
	Status   string            `json:"status"`
	UID      int64             `json:"uid"`
	Username string            `json:"username"` // user's first or preferred name
	Name     string            `json:"Name"`
	ImageURL string            `json:"ImageURL"`
	Message  string            `json:"message"`
	Token    string            `json:"Token"`
	Expire   util.JSONDateTime `json:"Expire"` // DATETIMEFMT in this format "2006-01-02T15:04:00Z"
}

// ValidateCookieResponse will be the response structure used when
// authentication is successful.
type ValidateCookieResponse struct {
	Status   string            `json:"status"`
	UID      int64             `json:"uid"`
	Name     string            `json:"Name"`
	ImageURL string            `json:"ImageURL"`
	Token    string            `json:"Token"`
	Expire   util.JSONDateTime `json:"Expire"` // DATETIMEFMT in this format "2006-01-02T15:04 "
}

// ValidateCookieRequest describes what the auth server wants to
// validate the cookie value
type ValidateCookieRequest struct {
	Status    string            `json:"status"`
	CookieVal string            `json:"cookieval"`
	IP        string            `json:"ip"`
	UserAgent string            `json:"useragent"`
	FLAGS     uint64            `json:"flags"`
	Expire    util.JSONDateTime `json:"expire"`
}
