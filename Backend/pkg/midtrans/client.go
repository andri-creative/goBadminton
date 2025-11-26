package midtrans

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Client interface {
	CreateTransaction(request *ChargeRequest) (*ChargeResponse, error)
	HandleNotification(payload map[string]interface{}) (*Notification, error)
}

type client struct {
	config     *Config
	httpClient *http.Client
}

func NewClient(config *Config) Client {
	return &client{
		config:     config,
		httpClient: &http.Client{},
	}
}

func (c *client) getBaseURL() string {
	if c.config.IsProduction() {
		return "https://api.midtrans.com/v2"
	}
	return "https://api.sandbox.midtrans.com/v2"
}

func (c *client) CreateTransaction(request *ChargeRequest) (*ChargeResponse, error) {
	url := c.getBaseURL() + "/charge"

	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.SetBasicAuth(c.config.ServerKey, "")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("midtrans API error: %s", string(body))
	}

	var chargeResponse ChargeResponse
	err = json.Unmarshal(body, &chargeResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	return &chargeResponse, nil
}

func (c *client) HandleNotification(payload map[string]interface{}) (*Notification, error) {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal notification: %v", err)
	}

	var notification Notification
	err = json.Unmarshal(jsonData, &notification)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal notification: %v", err)
	}

	return &notification, nil
}
