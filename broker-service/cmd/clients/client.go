package clients

import (
	"net/http"
	"time"
)

type Clients struct {
	Auth *AuthClient
	Log  *LogClient
}

func NewClients(authUrl, logUrl string) *Clients {

	httpClient := &http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 20,
			IdleConnTimeout:     90 * time.Second,
		},
	}

	return &Clients{
		Auth: NewAuthClient(authUrl, httpClient),
		Log:  NewLogClient(logUrl, httpClient),
	}
}
