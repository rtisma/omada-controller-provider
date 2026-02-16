package client

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"sync"
	"time"
)

const (
	defaultTimeout = 30 * time.Second
)

// Client represents an Omada Controller API client
type Client struct {
	baseURL      string
	username     string
	password     string
	siteID       string
	httpClient   *http.Client
	controllerID string
	csrfToken    string
	sessionID    string
	mu           sync.RWMutex
}

// ErrorResponse represents an API error response
type ErrorResponse struct {
	ErrorCode int    `json:"errorCode"`
	Message   string `json:"msg"`
}

// Error implements the error interface
func (e *ErrorResponse) Error() string {
	return fmt.Sprintf("API error %d: %s", e.ErrorCode, e.Message)
}

// APIResponse represents a standard API response wrapper
type APIResponse struct {
	ErrorCode int             `json:"errorCode"`
	Message   string          `json:"msg"`
	Result    json.RawMessage `json:"result,omitempty"`
}

// NewClient creates a new Omada Controller API client
func NewClient(host, username, password, siteID string, insecure bool) (*Client, error) {
	// Parse and validate the host URL
	baseURL, err := url.Parse(host)
	if err != nil {
		return nil, fmt.Errorf("invalid host URL: %w", err)
	}

	// Ensure HTTPS scheme
	if baseURL.Scheme == "" {
		baseURL.Scheme = "https"
	}

	// Create cookie jar for session management
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create cookie jar: %w", err)
	}

	// Create HTTP client with custom transport
	httpClient := &http.Client{
		Timeout: defaultTimeout,
		Jar:     jar,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: insecure,
			},
		},
	}

	client := &Client{
		baseURL:    baseURL.String(),
		username:   username,
		password:   password,
		siteID:     siteID,
		httpClient: httpClient,
	}

	// Authenticate and get controller ID
	if err := client.authenticate(); err != nil {
		return nil, fmt.Errorf("authentication failed: %w", err)
	}

	return client, nil
}

// doRequest performs an HTTP request with authentication
func (c *Client) doRequest(method, path string, body interface{}, result interface{}) error {
	c.mu.RLock()
	csrfToken := c.csrfToken
	c.mu.RUnlock()

	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(jsonBody)
	}

	reqURL := c.baseURL + path
	req, err := http.NewRequest(method, reqURL, bodyReader)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if csrfToken != "" {
		req.Header.Set("Csrf-Token", csrfToken)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	// Handle HTTP errors
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var errResp ErrorResponse
		if err := json.Unmarshal(respBody, &errResp); err == nil && errResp.ErrorCode != 0 {
			return &errResp
		}
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(respBody))
	}

	// Handle 401 Unauthorized - session expired
	if resp.StatusCode == 401 {
		c.mu.Lock()
		c.sessionID = ""
		c.csrfToken = ""
		c.mu.Unlock()

		// Try to re-authenticate
		if err := c.authenticate(); err != nil {
			return fmt.Errorf("re-authentication failed: %w", err)
		}

		// Retry the original request
		return c.doRequest(method, path, body, result)
	}

	// Parse the response
	if result != nil {
		var apiResp APIResponse
		if err := json.Unmarshal(respBody, &apiResp); err != nil {
			return fmt.Errorf("failed to parse response: %w", err)
		}

		if apiResp.ErrorCode != 0 {
			return &ErrorResponse{
				ErrorCode: apiResp.ErrorCode,
				Message:   apiResp.Message,
			}
		}

		if len(apiResp.Result) > 0 {
			if err := json.Unmarshal(apiResp.Result, result); err != nil {
				return fmt.Errorf("failed to parse result: %w", err)
			}
		}
	}

	return nil
}

// GetControllerID returns the controller ID
func (c *Client) GetControllerID() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.controllerID
}

// GetSiteID returns the site ID
func (c *Client) GetSiteID() string {
	return c.siteID
}
