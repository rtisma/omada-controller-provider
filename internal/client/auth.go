package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// InfoResponse represents the response from /api/info
type InfoResponse struct {
	ErrorCode    int    `json:"errorCode"`
	Message      string `json:"msg"`
	ControllerID string `json:"omadacId"`
}

// LoginRequest represents the login request payload
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse represents the login response
type LoginResponse struct {
	Token string `json:"token"`
}

// authenticate performs the full authentication flow
func (c *Client) authenticate() error {
	// Step 1: Get controller ID from /api/info
	if err := c.getControllerID(); err != nil {
		return fmt.Errorf("failed to get controller ID: %w", err)
	}

	// Step 2: Login to get session and CSRF token
	if err := c.login(); err != nil {
		return fmt.Errorf("login failed: %w", err)
	}

	return nil
}

// getControllerID retrieves the controller ID from /api/info
func (c *Client) getControllerID() error {
	reqURL := c.baseURL + "/api/info"

	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	var info InfoResponse
	if err := json.Unmarshal(body, &info); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	if info.ErrorCode != 0 {
		return fmt.Errorf("API error %d: %s", info.ErrorCode, info.Message)
	}

	if info.ControllerID == "" {
		return fmt.Errorf("controller ID not found in response")
	}

	c.mu.Lock()
	c.controllerID = info.ControllerID
	c.mu.Unlock()

	return nil
}

// login performs the login request to get session token and CSRF token
func (c *Client) login() error {
	c.mu.RLock()
	controllerID := c.controllerID
	c.mu.RUnlock()

	if controllerID == "" {
		return fmt.Errorf("controller ID not set")
	}

	loginReq := LoginRequest{
		Username: c.username,
		Password: c.password,
	}

	jsonBody, err := json.Marshal(loginReq)
	if err != nil {
		return fmt.Errorf("failed to marshal login request: %w", err)
	}

	reqURL := fmt.Sprintf("%s/%s/api/v2/login", c.baseURL, controllerID)

	req, err := http.NewRequest("POST", reqURL, bytes.NewReader(jsonBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	// Parse response to check for errors
	var apiResp APIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	if apiResp.ErrorCode != 0 {
		return fmt.Errorf("API error %d: %s", apiResp.ErrorCode, apiResp.Message)
	}

	// Extract CSRF token from response
	var loginResp LoginResponse
	if len(apiResp.Result) > 0 {
		if err := json.Unmarshal(apiResp.Result, &loginResp); err != nil {
			return fmt.Errorf("failed to parse login response: %w", err)
		}
	}

	// Store CSRF token
	c.mu.Lock()
	if loginResp.Token != "" {
		c.csrfToken = loginResp.Token
	}

	// Extract session cookie (TPOMADA_SESSIONID for 5.11+ or TPEAP_SESSIONID for older)
	cookies := resp.Cookies()
	for _, cookie := range cookies {
		if cookie.Name == "TPOMADA_SESSIONID" || cookie.Name == "TPEAP_SESSIONID" {
			c.sessionID = cookie.Value
			break
		}
	}
	c.mu.Unlock()

	if c.sessionID == "" {
		return fmt.Errorf("session cookie not found in login response")
	}

	return nil
}

// Logout terminates the current session
func (c *Client) Logout() error {
	c.mu.RLock()
	controllerID := c.controllerID
	c.mu.RUnlock()

	if controllerID == "" {
		return nil // Already logged out or not logged in
	}

	path := fmt.Sprintf("/%s/api/v2/logout", controllerID)

	// Logout doesn't need a response
	err := c.doRequest("POST", path, nil, nil)

	c.mu.Lock()
	c.sessionID = ""
	c.csrfToken = ""
	c.mu.Unlock()

	return err
}
