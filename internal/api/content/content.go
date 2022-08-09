package content

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/HUN220/magdam/internal/httpclient"
)

// `json:"content"` is returned when inline=true but is not currently
// supported
type ContentListItem struct {
	ContentId     string `json:"id"`
	ContentType   string `json:"type"`
	ContentLength int    `json:"length"`
}

type Result struct {
	Result string `json:"result"`
}

func GetContentList(c *httpclient.Client) ([]ContentListItem, error) {
	url := fmt.Sprintf("%s/api/v0/content/all", c.BaseURL)
	method := "GET"

	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	var res []ContentListItem

	if err := c.SendRequest(req, &res); err != nil {
		return nil, err
	}

	return res, nil
}

func GetContentBytes(c *httpclient.Client, id string, contentType string) ([]byte, error) {
	var content []byte

	var url string
	if strings.Contains(contentType, "json") {
		url = fmt.Sprintf("%s/api/v0/content/%s.json", c.BaseURL, id)
	} else if strings.Contains(contentType, "text/") {
		url = fmt.Sprintf("%s/api/v0/content/%s.text", c.BaseURL, id)
	} else {
		url = fmt.Sprintf("%s/api/v0/content/%s", c.BaseURL, id)
	}

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, err
	}

	// Add API Keys to request header
	req.Header.Add("X-Magda-API-Key-Id", c.ApiKeyId)
	req.Header.Add("X-Magda-API-Key", c.ApiKey)

	// GET Content from API
	res, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	// Try to unmarshall into Result
	if res.StatusCode != http.StatusOK {
		var errRes Result
		if err = json.NewDecoder(res.Body).Decode(&errRes); err == nil {
			return nil, errors.New(errRes.Result)
		}

		return nil, fmt.Errorf("unknown error, status code: %d", res.StatusCode)
	}

	// Read body of response
	content, _ = ioutil.ReadAll(res.Body)

	return content, nil
}

// Read the file at path and push it into Magda
func PutContentFile(c *httpclient.Client, id string, contentType string, path string) error {
	// Assemble API endpoint
	url := fmt.Sprintf("%s/api/v0/content/%s", c.BaseURL, id)

	// Load the local file
	data, err := os.Open(path)
	if err != nil {
		return err
	}
	defer data.Close()

	// Create request
	req, err := http.NewRequest("PUT", url, data)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", contentType)
	req.Header.Add("X-Magda-API-Key-Id", c.ApiKeyId)
	req.Header.Add("X-Magda-API-Key", c.ApiKey)

	// PUT into Magda
	res, err := c.HttpClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	// Try to unmarshall into Result
	if res.StatusCode != http.StatusCreated {
		var errRes Result
		if err = json.NewDecoder(res.Body).Decode(&errRes); err == nil {
			return errors.New(errRes.Result)
		}

		return fmt.Errorf("unknown error, status code: %d, id: %s", res.StatusCode, id)
	}

	return nil
}

// Return the full output filepath for a given ContentID and ContentType
func FilePath(id string, t string) string {
	switch t {
	case "application/json":
		return filepath.Join("data", fmt.Sprintf("%s.json", id))
	case "image/jpeg":
		return filepath.Join("data", fmt.Sprintf("%s.jpg", id))
	case "image/png":
		return filepath.Join("data", fmt.Sprintf("%s.png", id))
	case "image/x-icon":
		return filepath.Join("data", id)
	case "text/css":
		return filepath.Join("data", fmt.Sprintf("%s.css", id))
	case "text/html":
		return filepath.Join("data", fmt.Sprintf("%s.html", id))
	case "text/plain":
		return filepath.Join("data", fmt.Sprintf("%s.txt", id))
	default:
		return filepath.Join("data", id)
	}
}
