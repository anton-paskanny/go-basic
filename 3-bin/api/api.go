package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"demo/bin/config"
)

// Client defines interaction with JSON.BIN API
type Client interface {
	CreateBin(payload []byte, isPrivate bool) (string, error)
	GetBin(id string) ([]byte, error)
	UpdateBin(id string, payload []byte) error
	DeleteBin(id string) error
}

// client implements the Client interface
type client struct {
	httpClient *http.Client
	apiKey     string
}

const (
	baseURL          = "https://api.jsonbin.io/v3"
	headerMasterKey  = "X-Master-Key"
	headerBinPrivate = "X-Bin-Private"
)

// New creates a new API client with configuration
func New(cfg *config.Config) Client {
	return &client{
		httpClient: &http.Client{Timeout: 15 * time.Second},
		apiKey:     cfg.APIKey,
	}
}

// CreateBin sends a POST to create a new bin. Returns created bin id.
func (c *client) CreateBin(payload []byte, isPrivate bool) (string, error) {
	url := baseURL + "/b"
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(headerMasterKey, c.apiKey)
	if isPrivate {
		req.Header.Set(headerBinPrivate, "true")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
	}

	// JSONBin v3 returns {record: {...}, metadata: {id: "..."}}
	var decoded struct {
		Metadata struct {
			ID string `json:"id"`
		} `json:"metadata"`
	}
	if err := json.Unmarshal(body, &decoded); err != nil {
		return "", fmt.Errorf("decode response: %w", err)
	}
	if decoded.Metadata.ID == "" {
		return "", fmt.Errorf("no id in response")
	}
	return decoded.Metadata.ID, nil
}

// GetBin fetches bin contents as raw JSON
func (c *client) GetBin(id string) ([]byte, error) {
	url := baseURL + "/b/" + id
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set(headerMasterKey, c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
	}

	// JSONBin v3 returns {record: {...}, metadata: {...}}
	var decoded struct {
		Record json.RawMessage `json:"record"`
	}
	if err := json.Unmarshal(body, &decoded); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	if len(decoded.Record) == 0 {
		// Fallback to raw body if structure changes
		return body, nil
	}
	return decoded.Record, nil
}

// UpdateBin replaces bin contents with provided JSON
func (c *client) UpdateBin(id string, payload []byte) error {
	url := baseURL + "/b/" + id
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(headerMasterKey, c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
	}
	return nil
}

// DeleteBin removes a bin by id
func (c *client) DeleteBin(id string) error {
	url := baseURL + "/b/" + id
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	req.Header.Set(headerMasterKey, c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
	}
	return nil
}
