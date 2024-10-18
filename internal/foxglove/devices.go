package foxglove

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"
)

// ListDeviceResponse represents a robot device.
type ListDeviceResponse struct {
	ID         string            `json:"id"`
	Name       string            `json:"name"`
	OrgID      string            `json:"orgId"`
	CreatedAt  string            `json:"createdAt"`
	UpdatedAt  time.Time         `json:"updatedAt"`
	Properties map[string]string `json:"properties"`
}

// ListDevices fetches a list of devices with optional query parameters.
func (c *Client) ListDevices(query string, sortBy string, sortOrder string, limit int, offset int) ([]ListDeviceResponse, error) {
	params := url.Values{}
	if query != "" {
		params.Add("query", query)
	}
	if sortBy != "" {
		params.Add("sortBy", sortBy)
	}
	if sortOrder != "" {
		params.Add("sortOrder", sortOrder)
	}
	if limit > 0 {
		params.Add("limit", fmt.Sprintf("%d", limit))
	}
	if offset > 0 {
		params.Add("offset", fmt.Sprintf("%d", offset))
	}

	resp, err := c.doRequest("GET", "/devices?"+params.Encode(), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var devices []ListDeviceResponse
	if err := json.NewDecoder(resp.Body).Decode(&devices); err != nil {
		return nil, err
	}

	return devices, nil
}

// CreateDeviceRequest represents the payload to create a new device.
type CreateDeviceRequest struct {
	Name       string            `json:"name"`
	Properties map[string]string `json:"properties,omitempty"`
}

// CreateDeviceResponse represents the response returned after creating a new device.
type CreateDeviceResponse struct {
	ID         string            `json:"id"`
	Name       string            `json:"name"`
	OrgID      string            `json:"orgId"`
	CreatedAt  string            `json:"createdAt"`
	UpdatedAt  time.Time         `json:"updatedAt"`
	Properties map[string]string `json:"properties"`
}

// CreateDevice creates a new device with the specified name and properties.
func (c *Client) CreateDevice(reqBody CreateDeviceRequest) (*CreateDeviceResponse, error) {
	resp, err := c.doRequest("POST", "/devices", reqBody)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var createDeviceResp CreateDeviceResponse
	if err := json.NewDecoder(resp.Body).Decode(&createDeviceResp); err != nil {
		return nil, err
	}

	return &createDeviceResp, nil
}

// GetDeviceResponse represents the response returned when retrieving a specific device.
type GetDeviceResponse struct {
	ID         string            `json:"id"`
	Name       string            `json:"name"`
	OrgID      string            `json:"orgId"`
	CreatedAt  string            `json:"createdAt"`
	UpdatedAt  time.Time         `json:"updatedAt"`
	Properties map[string]string `json:"properties"`
}

// GetDevice retrieves the details of a specific device by its name or ID.
func (c *Client) GetDevice(nameOrId string) (*GetDeviceResponse, error) {
	encodedNameOrId := url.PathEscape(nameOrId)

	reqURL := fmt.Sprintf("/devices/%s", encodedNameOrId)

	resp, err := c.doRequest("GET", reqURL, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var getDeviceResp GetDeviceResponse
	if err := json.NewDecoder(resp.Body).Decode(&getDeviceResp); err != nil {
		return nil, err
	}

	return &getDeviceResp, nil
}

// UpdateDeviceRequest represents the payload to update a device.
type UpdateDeviceRequest struct {
	Name       string                 `json:"name,omitempty"`
	Properties map[string]interface{} `json:"properties,omitempty"`
}

// UpdateDeviceResponse represents the response returned when updating a device.
type UpdateDeviceResponse struct {
	ID         string                 `json:"id"`
	Name       string                 `json:"name"`
	OrgID      string                 `json:"orgId"`
	CreatedAt  string                 `json:"createdAt"`
	UpdatedAt  time.Time              `json:"updatedAt"`
	Properties map[string]interface{} `json:"properties"`
}

// UpdateDevice updates the details of a specific device by its name or ID.
func (c *Client) UpdateDevice(nameOrId string, reqBody UpdateDeviceRequest) (*UpdateDeviceResponse, error) {
	encodedNameOrId := url.PathEscape(nameOrId)

	reqURL := fmt.Sprintf("/devices/%s", encodedNameOrId)

	resp, err := c.doRequest("PATCH", reqURL, reqBody)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var updateDeviceResp UpdateDeviceResponse
	if err := json.NewDecoder(resp.Body).Decode(&updateDeviceResp); err != nil {
		return nil, err
	}

	return &updateDeviceResp, nil
}

// DeleteDeviceResponse represents the response returned when deleting a specific device.
type DeleteDeviceResponse struct {
	ID string `json:"id"`
}

// DeleteDevice deletes a device by its name or ID.
func (c *Client) DeleteDevice(nameOrId string) (*DeleteDeviceResponse, error) {
	encodedNameOrId := url.PathEscape(nameOrId)

	reqURL := fmt.Sprintf("/devices/%s", encodedNameOrId)

	resp, err := c.doRequest("DELETE", reqURL, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var deleteDeviceResp DeleteDeviceResponse
	if err := json.NewDecoder(resp.Body).Decode(&deleteDeviceResp); err != nil {
		return nil, err
	}

	return &deleteDeviceResp, nil
}
