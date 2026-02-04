package clients

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type AuthClient struct {
	baseUrl string
	client  *http.Client
}

func NewAuthClient(url string, client *http.Client) *AuthClient {

	return &AuthClient{
		baseUrl: url,
		client:  client,
	}
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func (auth *AuthClient) Login(ctx context.Context, payload *AuthPayload) (*AuthResponse, error) {
	url := auth.baseUrl + "/authenticate"

	// Json request payload for auth-service
	jsonData, err := json.Marshal(*payload)

	if err != nil {
		return nil, err
	}
	fmt.Println(jsonData, url)

	// Create a request with context
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonData))

	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Send Request using client
	resp, err := auth.client.Do(req)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check if the status is ok
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Auth service return status code %d", resp.StatusCode)
	}

	var result AuthResponse

	json.NewDecoder(resp.Body).Decode(&result)

	return &result, nil

}
