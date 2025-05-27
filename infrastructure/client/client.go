package client

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"messager/domain/message"
)

type Client interface {
	SendMessage(ctx context.Context, message message.Message) (id string, err error)
}

type Config struct {
	URL     string
	Token   string
	Timeout time.Duration
}

type client struct {
	client *http.Client
	config *Config
}

type requestPayload struct {
	To      string `json:"to"`
	Content string `json:"content"`
}

type responsePayload struct {
	MessageID string `json:"messageId"`
}

func New(config Config) Client {
	return &client{
		client: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
			Timeout: config.Timeout,
		},
		config: &config,
	}
}

func (c *client) SendMessage(ctx context.Context, message message.Message) (id string, err error) {
	body, err := json.Marshal(requestPayload{
		To:      message.Phone,
		Content: message.Content,
	})
	if err != nil {
		return "", fmt.Errorf("json.Marshal(): %w", err)
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, c.config.URL, bytes.NewBuffer(body))
	if err != nil {
		return "", fmt.Errorf("http.NewRequestWithContext(): %w", err)
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("x-ins-auth-key", c.config.Token)

	response, err := c.client.Do(request)
	if err != nil {
		return "", fmt.Errorf("http.Client.Do(): %w", err)
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		return "", fmt.Errorf("http.Client.Do(): unexpected response status code %d", response.StatusCode)
	}

	var payload responsePayload

	if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
		return "", fmt.Errorf("json.NewDecoder().Decode(): %w", err)
	}

	return payload.MessageID, nil
}
