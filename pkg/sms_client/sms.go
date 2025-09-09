package smsclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type SmsClient struct {
	Host       string
	HttpClient *http.Client
}

func NewSmsClient(host string) *SmsClient {
	return &SmsClient{
		Host: host,
		HttpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *SmsClient) SendSMS(to string, message string) (string, error) {
	payload := map[string]string{
		"to":      to,
		"content": message,
	}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest("POST", c.Host, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Since I couldn't get response like given in PDF document. I wanted to simulate it via using x-request-id header.
	messageId := resp.Header.Get("x-request-id")

	return messageId, nil
}
