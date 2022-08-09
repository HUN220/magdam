package pull

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/HUN220/magdam/internal/api/content"
	"github.com/HUN220/magdam/internal/httpclient"
)

// Handle "pull" subcommand
func PullCmd(apiKeyIdPtr string, apiKeyPtr string, baseUrlPtr string) {
	c := httpclient.NewClient(apiKeyIdPtr, apiKeyPtr, baseUrlPtr)

	// Pull list of Content from Magda and write to disk
	res, err := pullContentList(c)
	if err != nil {
		log.Fatal(err)
	}

	// Pull each item from Content list
	for _, item := range res {
		err := pullContentItem(c, item)
		if err != nil {
			log.Fatal(err)
		}
	}
}

// Pull the Content item from the API and write to the output directory
func pullContentItem(c *httpclient.Client, item content.ContentListItem) error {
	if item.ContentLength <= 0 {
		return nil
	}

	// Get content data
	res, err := content.GetContentBytes(c, item.ContentId, item.ContentType)
	if err != nil {
		log.Println(fmt.Sprintf("%s: Unable to get Content from API (%s)", item.ContentId, err))
		return nil
	}

	fullPath := content.FilePath(item.ContentId, item.ContentType)

	// Create directory structure for content
	if strings.Contains(item.ContentId, "/") {
		dataDir := filepath.Dir(fullPath)
		_, err := os.Stat(dataDir)
		if os.IsNotExist(err) {
			errDir := os.MkdirAll(dataDir, 0755)
			if errDir != nil {
				log.Fatal(fmt.Sprintf("%s: Unable to create content directory (%s)", item.ContentId, err))
			}
		}
	}

	// Write data to new directory
	err = ioutil.WriteFile(fullPath, res, 0644)
	if err != nil {
		log.Fatal(fmt.Sprintf("%s: Unable to write Content to file (%s)", item.ContentId, err))
	}

	return nil
}

// Pull the list of Content from the API and write to the output directory
func pullContentList(c *httpclient.Client) ([]content.ContentListItem, error) {
	// Get list of Content from Magda
	res, _ := content.GetContentList(c)

	// Create data output directory
	_, err := os.Stat("data")
	if os.IsNotExist(err) {
		errDir := os.MkdirAll("data", 0755)
		if errDir != nil {
			return nil, fmt.Errorf("failed to create output directory 'data' (%s)", err)
		}
	}

	// Write out list of content to json
	contFile, err := json.MarshalIndent(res, "", " ")
	if err != nil {
		return nil, fmt.Errorf("failed to mashal API response as JSON (%s)", err)
	}

	output := "./data/content_all.json"

	err = ioutil.WriteFile(output, contFile, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to write to '%s' (%s)", output, err)
	}

	return res, nil
}
