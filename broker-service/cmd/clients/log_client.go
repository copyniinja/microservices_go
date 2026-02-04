package clients

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type LogClient struct {
	baseUrl string
	client  *http.Client
}

type LogPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}
type LogResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
}

func NewLogClient(url string, client *http.Client) *LogClient {

	return &LogClient{
		baseUrl: url,
		client:  client,
	}
}

func (l *LogClient) Insert(ctx context.Context, payload *LogPayload) (*LogResponse, error) {

	url := l.baseUrl + "/log"

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
	resp, err := l.client.Do(req)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		return nil, fmt.Errorf("Logger service return status code %d", resp.StatusCode)
	}

	var result LogResponse

	json.NewDecoder(resp.Body).Decode(&result)

	return &result, nil

}
