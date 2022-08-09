package httpclient

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

// Client .
type Client struct {
	ApiKeyId   string
	ApiKey     string
	BaseURL    string
	HttpClient *http.Client
}

// NewClient creates new Magda client with given API key
func NewClient(apiKeyId string, apiKey string, baseURL string) *Client {
	return &Client{
		ApiKeyId: apiKeyId,
		ApiKey:   apiKey,
		BaseURL:  baseURL,
		HttpClient: &http.Client{
			Timeout: 2 * time.Minute,
		},
	}
}

type errorResponse struct {
	Result string `json:"result"`
}

// Content-type and body should be already added to req
func (c *Client) SendRequest(req *http.Request, v interface{}) error {
	req.Header.Set("Accept", "application/json; charset=utf-8")
	req.Header.Add("X-Magda-API-Key-Id", c.ApiKeyId)
	req.Header.Add("X-Magda-API-Key", c.ApiKey)

	res, err := c.HttpClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	// Try to unmarshall into errorResponse
	if res.StatusCode != http.StatusOK {
		var errRes errorResponse
		if err = json.NewDecoder(res.Body).Decode(&errRes); err == nil {
			return errors.New(errRes.Result)
		}

		return fmt.Errorf("unknown error, status code: %d", res.StatusCode)
	}

	if err = json.NewDecoder(res.Body).Decode(&v); err != nil {
		return err
	}

	return nil
}
