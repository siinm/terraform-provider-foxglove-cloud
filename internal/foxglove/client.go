package foxglove

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"strings"
)

// Client represents the API client.
type Client struct {
	BaseURL string
	APIKey  string
	Client  *http.Client
}

type loggingTransport struct{}

func (s *loggingTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	bytes, _ := httputil.DumpRequestOut(r, true)

	resp, err := http.DefaultTransport.RoundTrip(r)
	// err is returned after dumping the response

	respBytes, _ := httputil.DumpResponse(resp, true)
	bytes = append(bytes, respBytes...)

	fmt.Printf("%s\n", bytes)

	return resp, err
}

// NewClient initializes and returns a new API client.
func NewClient(apiKey string) *Client {
	return &Client{
		BaseURL: "https://api.foxglove.dev/v1", // Replace with actual base URL
		APIKey:  apiKey,
		Client: &http.Client{
			Transport: &loggingTransport{},
		},
	}
}

func (c *Client) doRequest(method string, uri string, reqBody interface{}) (*http.Response, error) {
	url := c.BaseURL + uri
	var body io.Reader = nil
	if reqBody != nil {
		reqBodyJSON, err := json.Marshal(reqBody)
		if err != nil {
			return nil, err
		}
		body = strings.NewReader(string(reqBodyJSON))
	}
	req, err := http.NewRequest(method, url, body)

	if err != nil {
		return nil, err
	}

	if strings.HasPrefix(c.APIKey, "fox.session=") {
		// this is for when you authenticate using a session cookie
		req.AddCookie(&http.Cookie{
			Name:  "fox.session",
			Value: strings.TrimPrefix(c.APIKey, "fox.session="),
		})
	} else {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.APIKey))
	}

	if reqBody != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		defer resp.Body.Close()
		respBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to send request to %s: status code %d %s", url, resp.StatusCode, string(respBytes))
	}
	return resp, nil
}
