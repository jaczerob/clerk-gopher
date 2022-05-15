package login

import (
	"net/http"
	"net/url"
)

type LoginClient struct {
	http    *http.Client
	headers map[string]string
	baseURL *url.URL
}

type LoginData struct {
	Success    string `json:"success,omitempty"`
	Gameserver string `json:"gameserver,omitempty"`
	Playcookie string `json:"cookie,omitempty"`
	AppToken   string `json:"appToken,omitempty"`
	AuthToken  string `json:"authToken,omitempty"`
	ETA        string `json:"eta,omitempty"`
	Position   string `json:"position,omitempty"`
	QueueToken string `json:"queueToken,omitempty"`
	Banner     string `json:"banner,omitempty"`
}
