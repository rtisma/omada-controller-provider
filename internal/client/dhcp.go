package client

import (
	"fmt"
)

// DHCPReservation represents a DHCP reservation (static IP assignment)
type DHCPReservation struct {
	ID         string `json:"id,omitempty"`
	Name       string `json:"name"`
	MACAddress string `json:"mac"`
	IPAddress  string `json:"ip"`
	NetworkID  string `json:"networkId"`
	Comment    string `json:"comment,omitempty"`
}

// DHCPReservationsResponse represents the response from the DHCP reservations list endpoint
type DHCPReservationsResponse struct {
	Data []DHCPReservation `json:"data"`
}

// CreateDHCPReservation creates a new DHCP reservation
func (c *Client) CreateDHCPReservation(siteID string, reservation *DHCPReservation) (*DHCPReservation, error) {
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

	path := fmt.Sprintf("/%s/api/v2/sites/%s/setting/lan/ipReservations", controllerID, siteID)

	var result DHCPReservation
	if err := c.doRequest("POST", path, reservation, &result); err != nil {
		return nil, fmt.Errorf("failed to create DHCP reservation: %w", err)
	}

	return &result, nil
}

// GetDHCPReservation retrieves a DHCP reservation by ID
func (c *Client) GetDHCPReservation(siteID, reservationID string) (*DHCPReservation, error) {
	reservations, err := c.GetDHCPReservations(siteID)
	if err != nil {
		return nil, err
	}

	for _, reservation := range reservations {
		if reservation.ID == reservationID {
			return &reservation, nil
		}
	}

	return nil, fmt.Errorf("DHCP reservation not found: %s", reservationID)
}

// GetDHCPReservations retrieves all DHCP reservations for a site
func (c *Client) GetDHCPReservations(siteID string) ([]DHCPReservation, error) {
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

	path := fmt.Sprintf("/%s/api/v2/sites/%s/setting/lan/ipReservations", controllerID, siteID)

	var response DHCPReservationsResponse
	if err := c.doRequest("GET", path, nil, &response); err != nil {
		return nil, fmt.Errorf("failed to get DHCP reservations: %w", err)
	}

	return response.Data, nil
}

// UpdateDHCPReservation updates an existing DHCP reservation
func (c *Client) UpdateDHCPReservation(siteID, reservationID string, reservation *DHCPReservation) (*DHCPReservation, error) {
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

	path := fmt.Sprintf("/%s/api/v2/sites/%s/setting/lan/ipReservations/%s", controllerID, siteID, reservationID)

	var result DHCPReservation
	if err := c.doRequest("PATCH", path, reservation, &result); err != nil {
		return nil, fmt.Errorf("failed to update DHCP reservation: %w", err)
	}

	return &result, nil
}

// DeleteDHCPReservation deletes a DHCP reservation
func (c *Client) DeleteDHCPReservation(siteID, reservationID string) error {
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

	path := fmt.Sprintf("/%s/api/v2/sites/%s/setting/lan/ipReservations/%s", controllerID, siteID, reservationID)

	if err := c.doRequest("DELETE", path, nil, nil); err != nil {
		return fmt.Errorf("failed to delete DHCP reservation: %w", err)
	}

	return nil
}
