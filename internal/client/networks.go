package client

import (
	"fmt"
)

// Network represents an Omada network/VLAN
type Network struct {
	ID           string `json:"id,omitempty"`
	Name         string `json:"name"`
	VlanID       int    `json:"vlan"`
	Gateway      string `json:"gateway"`
	Netmask      string `json:"netmask"`
	DHCPEnabled  bool   `json:"dhcpEnable"`
	DHCPStart    string `json:"dhcpStart,omitempty"`
	DHCPEnd      string `json:"dhcpEnd,omitempty"`
	LeaseTime    int    `json:"leaseTime,omitempty"`
	DNSPrimary   string `json:"primaryDns,omitempty"`
	DNSSecondary string `json:"secondaryDns,omitempty"`
	DomainName   string `json:"domainName,omitempty"`
	Purpose      string `json:"purpose,omitempty"`
	NetworkType  string `json:"networkType,omitempty"`
}

// NetworksResponse represents the response from the networks list endpoint
type NetworksResponse struct {
	Data []Network `json:"data"`
}

// CreateNetwork creates a new network/VLAN
func (c *Client) CreateNetwork(siteID string, network *Network) (*Network, error) {
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

	path := fmt.Sprintf("/%s/api/v2/sites/%s/setting/lan/networks", controllerID, siteID)

	var result Network
	if err := c.doRequest("POST", path, network, &result); err != nil {
		return nil, fmt.Errorf("failed to create network: %w", err)
	}

	return &result, nil
}

// GetNetwork retrieves a network by ID
func (c *Client) GetNetwork(siteID, networkID string) (*Network, error) {
	networks, err := c.GetNetworks(siteID)
	if err != nil {
		return nil, err
	}

	for _, network := range networks {
		if network.ID == networkID {
			return &network, nil
		}
	}

	return nil, fmt.Errorf("network not found: %s", networkID)
}

// GetNetworks retrieves all networks for a site
func (c *Client) GetNetworks(siteID string) ([]Network, error) {
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

	path := fmt.Sprintf("/%s/api/v2/sites/%s/setting/lan/networks", controllerID, siteID)

	var response NetworksResponse
	if err := c.doRequest("GET", path, nil, &response); err != nil {
		return nil, fmt.Errorf("failed to get networks: %w", err)
	}

	return response.Data, nil
}

// UpdateNetwork updates an existing network
func (c *Client) UpdateNetwork(siteID, networkID string, network *Network) (*Network, error) {
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

	path := fmt.Sprintf("/%s/api/v2/sites/%s/setting/lan/networks/%s", controllerID, siteID, networkID)

	var result Network
	if err := c.doRequest("PATCH", path, network, &result); err != nil {
		return nil, fmt.Errorf("failed to update network: %w", err)
	}

	return &result, nil
}

// DeleteNetwork deletes a network
func (c *Client) DeleteNetwork(siteID, networkID string) error {
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

	path := fmt.Sprintf("/%s/api/v2/sites/%s/setting/lan/networks/%s", controllerID, siteID, networkID)

	if err := c.doRequest("DELETE", path, nil, nil); err != nil {
		return fmt.Errorf("failed to delete network: %w", err)
	}

	return nil
}
