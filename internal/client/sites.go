package client

import (
	"fmt"
)

// Site represents an Omada site
type Site struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	Location string `json:"location"`
	TimeZone string `json:"timezone"`
	Scenario string `json:"scenario"`
}

// SitesResponse represents the response from the sites list endpoint
type SitesResponse struct {
	Data []Site `json:"data"`
}

// GetSites retrieves all sites
func (c *Client) GetSites() ([]Site, error) {
	c.mu.RLock()
	controllerID := c.controllerID
	c.mu.RUnlock()

	if controllerID == "" {
		return nil, fmt.Errorf("not authenticated")
	}

	path := fmt.Sprintf("/%s/api/v2/sites", controllerID)

	var response SitesResponse
	if err := c.doRequest("GET", path, nil, &response); err != nil {
		return nil, fmt.Errorf("failed to get sites: %w", err)
	}

	return response.Data, nil
}

// GetSite retrieves a specific site by name
func (c *Client) GetSite(siteName string) (*Site, error) {
	sites, err := c.GetSites()
	if err != nil {
		return nil, err
	}

	for _, site := range sites {
		if site.Name == siteName {
			return &site, nil
		}
	}

	return nil, fmt.Errorf("site not found: %s", siteName)
}
