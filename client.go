package audittrail

import (
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
)

// AuditTrailClient is a client for interacting with the AuditTrail API.
type AuditTrailClient struct {
	client  *resty.Client
	baseURL string
	token   string
}

// LogListResponse represents the paginated response from the list logs endpoint.
type LogListResponse struct {
	Count    int         `json:"count"`
	Next     *string     `json:"next"`
	Previous *string     `json:"previous"`
	Results  []LogEntry  `json:"results"`
}

// LogEntry represents a single log entry in the audit trail.
type LogEntry struct {
	ID        string                 `json:"_id"`
	Timestamp string                 `json:"timestamp"`
	Action    string                 `json:"action"`
	UserID    string                 `json:"user_id"`
	Details   map[string]interface{} `json:"details"`
	Signature string                 `json:"signature"`
}

// NewAuditTrailClient initializes a new AuditTrailClient with the given base URL and JWT token.
// The baseURL must be a valid URL (e.g., "https://api.example.com" or "http://localhost:8080")
// where the AuditTrail API is running.
func NewAuditTrailClient(baseURL, token string) (*AuditTrailClient, error) {
	if baseURL == "" {
		return nil, errors.New("baseURL is required")
	}
	if _, err := url.ParseRequestURI(baseURL); err != nil {
		return nil, fmt.Errorf("invalid baseURL: %w", err)
	}

	client := resty.New()
	client.
		SetRetryCount(3).
		SetRetryWaitTime(1 * time.Second).
		SetRetryMaxWaitTime(10 * time.Second).
		SetBaseURL(baseURL).
		SetHeader("Authorization", fmt.Sprintf("Bearer %s", token)).
		SetHeader("Content-Type", "application/json")
	return &AuditTrailClient{
		client:  client,
		baseURL: baseURL,
		token:   token,
	}, nil
}

// CreateLog creates a new log entry in the audit trail.
func (c *AuditTrailClient) CreateLog(action string, details map[string]interface{}) (map[string]string, error) {
	if c.baseURL == "" {
		return nil, errors.New("client not properly initialized: baseURL is missing")
	}
	body := map[string]interface{}{
		"action": action,
	}
	if details != nil {
		body["details"] = details
	}

	requestID := uuid.New().String()
	resp, err := c.client.R().
		SetHeader("X-Request-ID", requestID).
		SetBody(body).
		SetResult(&map[string]string{}).
		Post("/api/logs/")
	if err != nil {
		return nil, fmt.Errorf("failed to create log: %w", err)
	}
	return *resp.Result().(*map[string]string), nil
}

// ListLogs retrieves a list of logs with optional filters.
func (c *AuditTrailClient) ListLogs(params map[string]string) (*LogListResponse, error) {
	if c.baseURL == "" {
		return nil, errors.New("client not properly initialized: baseURL is missing")
	}
	requestID := uuid.New().String()
	resp, err := c.client.R().
		SetHeader("X-Request-ID", requestID).
		SetQueryParams(params).
		SetResult(&LogListResponse{}).
		Get("/api/logs/list/")
	if err != nil {
		return nil, fmt.Errorf("failed to list logs: %w", err)
	}
	return resp.Result().(*LogListResponse), nil
}