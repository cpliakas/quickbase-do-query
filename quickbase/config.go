package quickbase

import "net/http"

// Config contains the configuration used to make API requests.
type Config struct {
	HTTPClient *http.Client
	AppToken   string
	Ticket     string
	UserToken  string
	RealmHost  string
	AppID      string
}
