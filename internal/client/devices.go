package client

import (
	"fmt"
)

// Device represents an Omada device (AP, switch, gateway)
type Device struct {
	MAC             string `json:"mac"`
	Name            string `json:"name"`
	Type            string `json:"type"`
	Model           string `json:"model"`
	Status          string `json:"status"`
	LEDEnabled      bool   `json:"ledSetting"`
	Location        string `json:"location,omitempty"`
	Site            string `json:"site"`
	IP              string `json:"ip,omitempty"`
	Uptime          int64  `json:"uptime,omitempty"`
	FirmwareVersion string `json:"firmwareVersion,omitempty"`
	Adoption        bool   `json:"needAdopt"`
}

// DevicesResponse represents the response from the devices list endpoint
type DevicesResponse struct {
	Data []Device `json:"data"`
}

// DeviceConfig represents configuration for a device
type DeviceConfig struct {
	Name       string `json:"name"`
	LEDEnabled bool   `json:"ledSetting"`
	Location   string `json:"location,omitempty"`
}

// GetDevices retrieves all devices for a site
func (c *Client) GetDevices(siteID string) ([]Device, error) {
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

	path := fmt.Sprintf("/%s/api/v2/sites/%s/devices", controllerID, siteID)

	var response DevicesResponse
	if err := c.doRequest("GET", path, nil, &response); err != nil {
		return nil, fmt.Errorf("failed to get devices: %w", err)
	}

	return response.Data, nil
}

// GetDevice retrieves a specific device by MAC address
func (c *Client) GetDevice(siteID, mac string) (*Device, error) {
	devices, err := c.GetDevices(siteID)
	if err != nil {
		return nil, err
	}

	for _, device := range devices {
		if device.MAC == mac {
			return &device, nil
		}
	}

	return nil, fmt.Errorf("device not found: %s", mac)
}

// UpdateDevice updates device configuration
func (c *Client) UpdateDevice(siteID, mac string, config *DeviceConfig) (*Device, error) {
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

	path := fmt.Sprintf("/%s/api/v2/sites/%s/devices/%s", controllerID, siteID, mac)

	var result Device
	if err := c.doRequest("PATCH", path, config, &result); err != nil {
		return nil, fmt.Errorf("failed to update device: %w", err)
	}

	return &result, nil
}

// AdoptDevice adopts a pending device
func (c *Client) AdoptDevice(siteID, mac string) error {
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

	path := fmt.Sprintf("/%s/api/v2/sites/%s/cmd/devices/%s/adopt", controllerID, siteID, mac)

	if err := c.doRequest("POST", path, nil, nil); err != nil {
		return fmt.Errorf("failed to adopt device: %w", err)
	}

	return nil
}

// ForgetDevice forgets (removes) a device
func (c *Client) ForgetDevice(siteID, mac string) error {
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

	path := fmt.Sprintf("/%s/api/v2/sites/%s/cmd/devices/%s/forget", controllerID, siteID, mac)

	if err := c.doRequest("POST", path, nil, nil); err != nil {
		return fmt.Errorf("failed to forget device: %w", err)
	}

	return nil
}

// RebootDevice reboots a device
func (c *Client) RebootDevice(siteID, mac string) error {
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

	path := fmt.Sprintf("/%s/api/v2/sites/%s/cmd/devices/%s/reboot", controllerID, siteID, mac)

	if err := c.doRequest("POST", path, nil, nil); err != nil {
		return fmt.Errorf("failed to reboot device: %w", err)
	}

	return nil
}
