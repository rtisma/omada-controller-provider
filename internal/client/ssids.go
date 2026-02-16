package client

import (
	"fmt"
)

// SSID represents an Omada wireless network
type SSID struct {
	ID              string `json:"id,omitempty"`
	Name            string `json:"name"`
	SSID            string `json:"ssid"`
	Enabled         bool   `json:"enable"`
	HideSSID        bool   `json:"hideSSID"`
	SecurityMode    string `json:"securityMode"`
	Password        string `json:"password,omitempty"`
	VlanID          int    `json:"vlanId,omitempty"`
	GuestNetwork    bool   `json:"guestNetwork"`
	ClientIsolation bool   `json:"clientIsolation"`
	Band2_4GEnabled bool   `json:"band24gEnable"`
	Band5GEnabled   bool   `json:"band5gEnable"`
	Band6GEnabled   bool   `json:"band6gEnable"`
	MaxClients      int    `json:"maxClients,omitempty"`
	RateLimit       bool   `json:"rateLimit"`
	DownlinkLimit   int    `json:"downlinkLimit,omitempty"`
	UplinkLimit     int    `json:"uplinkLimit,omitempty"`
	ScheduleEnabled bool   `json:"scheduleEnable"`
	PortalEnabled   bool   `json:"portalEnable"`
	RadiusProfile   string `json:"radiusProfile,omitempty"`
}

// SSIDsResponse represents the response from the SSIDs list endpoint
type SSIDsResponse struct {
	Data []SSID `json:"data"`
}

// CreateSSID creates a new wireless network
func (c *Client) CreateSSID(siteID string, ssid *SSID) (*SSID, error) {
	c.mu.RLock()
	controllerID := c.controllerID
	c.mu.RUnlock()

	if controllerID == "" {
		return nil, fmt.Errorf("not authenticated")
	}

	// Use the site ID from the client if not provided
	if siteID == "" {
		siteID = c.siteID
	}

	path := fmt.Sprintf("/%s/api/v2/sites/%s/setting/wlans", controllerID, siteID)

	var result SSID
	if err := c.doRequest("POST", path, ssid, &result); err != nil {
		return nil, fmt.Errorf("failed to create SSID: %w", err)
	}

	return &result, nil
}

// GetSSID retrieves an SSID by ID
func (c *Client) GetSSID(siteID, ssidID string) (*SSID, error) {
	ssids, err := c.GetSSIDs(siteID)
	if err != nil {
		return nil, err
	}

	for _, ssid := range ssids {
		if ssid.ID == ssidID {
			return &ssid, nil
		}
	}

	return nil, fmt.Errorf("SSID not found: %s", ssidID)
}

// GetSSIDs retrieves all SSIDs for a site
func (c *Client) GetSSIDs(siteID string) ([]SSID, error) {
	c.mu.RLock()
	controllerID := c.controllerID
	c.mu.RUnlock()

	if controllerID == "" {
		return nil, fmt.Errorf("not authenticated")
	}

	// Use the site ID from the client if not provided
	if siteID == "" {
		siteID = c.siteID
	}

	path := fmt.Sprintf("/%s/api/v2/sites/%s/setting/wlans", controllerID, siteID)

	var response SSIDsResponse
	if err := c.doRequest("GET", path, nil, &response); err != nil {
		return nil, fmt.Errorf("failed to get SSIDs: %w", err)
	}

	return response.Data, nil
}

// UpdateSSID updates an existing SSID
func (c *Client) UpdateSSID(siteID, ssidID string, ssid *SSID) (*SSID, error) {
	c.mu.RLock()
	controllerID := c.controllerID
	c.mu.RUnlock()

	if controllerID == "" {
		return nil, fmt.Errorf("not authenticated")
	}

	// Use the site ID from the client if not provided
	if siteID == "" {
		siteID = c.siteID
	}

	path := fmt.Sprintf("/%s/api/v2/sites/%s/setting/wlans/%s", controllerID, siteID, ssidID)

	var result SSID
	if err := c.doRequest("PATCH", path, ssid, &result); err != nil {
		return nil, fmt.Errorf("failed to update SSID: %w", err)
	}

	return &result, nil
}

// DeleteSSID deletes an SSID
func (c *Client) DeleteSSID(siteID, ssidID string) error {
	c.mu.RLock()
	controllerID := c.controllerID
	c.mu.RUnlock()

	if controllerID == "" {
		return fmt.Errorf("not authenticated")
	}

	// Use the site ID from the client if not provided
	if siteID == "" {
		siteID = c.siteID
	}

	path := fmt.Sprintf("/%s/api/v2/sites/%s/setting/wlans/%s", controllerID, siteID, ssidID)

	if err := c.doRequest("DELETE", path, nil, nil); err != nil {
		return fmt.Errorf("failed to delete SSID: %w", err)
	}

	return nil
}
