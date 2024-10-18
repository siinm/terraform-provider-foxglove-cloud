package foxglove

import (
	"encoding/json"
	"fmt"
	"net/url"
)

// ListAPIKeyResponse represents an API key.
type ListAPIKeyResponse struct {
	ID                   string   `json:"id"`
	OrgID                string   `json:"orgId"`
	Label                string   `json:"label"`
	Capabilities         []string `json:"capabilities"`
	CreatedAt            string   `json:"createdAt"`
	UpdatedAt            string   `json:"updatedAt"`
	LastSeenAt           string   `json:"lastSeenAt"`
	Enabled              bool     `json:"enabled"`
	CreatedByOrgMemberId string   `json:"createdByOrgMemberId"`
}

// ListAPIKeys fetches a list of API keys.
func (c *Client) ListAPIKeys() ([]ListAPIKeyResponse, error) {
	resp, err := c.doRequest("GET", "/api-keys", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var apiKeys []ListAPIKeyResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiKeys); err != nil {
		return nil, err
	}

	return apiKeys, nil
}

// CreateAPIKeyRequest represents the payload to create a new API key.
type CreateAPIKeyRequest struct {
	Label        string   `json:"label"`
	Capabilities []string `json:"capabilities"`
}

// CreateAPIKeyResponse represents the response returned after creating a new API key.
type CreateAPIKeyResponse struct {
	ID                   string   `json:"id"`
	OrgID                string   `json:"orgId"`
	Label                string   `json:"label"`
	Capabilities         []string `json:"capabilities"`
	CreatedAt            string   `json:"createdAt"`
	UpdatedAt            string   `json:"updatedAt"`
	CreatedByOrgMemberId string   `json:"createdByOrgMemberId"`
	SecretToken          string   `json:"secretToken"`
}

// CreateAPIKey creates a new API key with the specified label and capabilities.
func (c *Client) CreateAPIKey(reqBody CreateAPIKeyRequest) (*CreateAPIKeyResponse, error) {
	resp, err := c.doRequest("POST", "/api-keys", reqBody)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var createAPIKeyResp CreateAPIKeyResponse
	if err := json.NewDecoder(resp.Body).Decode(&createAPIKeyResp); err != nil {
		return nil, err
	}

	return &createAPIKeyResp, nil
}

// UpdateAPIKeyRequest represents the payload to update an API key.
type UpdateAPIKeyRequest struct {
	Label        string   `json:"label,omitempty"`
	Capabilities []string `json:"capabilities,omitempty"`
}

// UpdateAPIKeyResponse represents the response returned when updating an API key.
type UpdateAPIKeyResponse struct {
	ID                   string   `json:"id"`
	OrgID                string   `json:"orgId"`
	Label                string   `json:"label"`
	Capabilities         []string `json:"capabilities"`
	CreatedAt            string   `json:"createdAt"`
	UpdatedAt            string   `json:"updatedAt"`
	Enabled              bool     `json:"enabled"`
	CreatedByOrgMemberId string   `json:"createdByOrgMemberId"`
}

// UpdateAPIKey updates the details of a specific API key by its ID.
func (c *Client) UpdateAPIKey(id string, reqBody UpdateAPIKeyRequest) (*UpdateAPIKeyResponse, error) {
	encodedID := url.PathEscape(id)

	reqURL := fmt.Sprintf("/api-keys/%s", encodedID)

	resp, err := c.doRequest("PATCH", reqURL, reqBody)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var updateAPIKeyResp UpdateAPIKeyResponse
	if err := json.NewDecoder(resp.Body).Decode(&updateAPIKeyResp); err != nil {
		return nil, err
	}

	return &updateAPIKeyResp, nil
}

// DeleteAPIKey deletes an API key by its ID.
func (c *Client) DeleteAPIKey(id string) error {
	encodedID := url.PathEscape(id)

	reqURL := fmt.Sprintf("/api-keys/%s", encodedID)

	resp, err := c.doRequest("DELETE", reqURL, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// API doesn't return any content on successful deletion, so just check for errors
	return nil
}
