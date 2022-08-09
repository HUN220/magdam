package push

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/HUN220/magdam/internal/api/content"
	"github.com/HUN220/magdam/internal/httpclient"
)

// Handle "push" subcommand
func PushCmd(apiKeyId string, apiKey string, baseUrl string) {
	c := httpclient.NewClient(apiKeyId, apiKey, baseUrl)

	// Read the Content list from the `content_all.json`
	var content []content.ContentListItem
	err := loadContentList(&content)
	if err != nil {
		log.Fatal(err)
	}

	// Push each item from the data directory into Magda
	for _, item := range content {
		err := pushContentItem(c, item)
		if err != nil {
			log.Fatal(err)
		}
	}
}

// Read the Content file and push it into Magda via the API
func pushContentItem(c *httpclient.Client, item content.ContentListItem) error {
	// Magda doesn't appear to support an empty payload
	if item.ContentLength <= 0 {
		return nil
	}

	path := content.FilePath(item.ContentId, item.ContentType)
	err := content.PutContentFile(c, item.ContentId, item.ContentType, path)

	if err != nil {
		return err
	}

	return nil
}

// Populate []content.ContentListItem using values from ./data/content_all.json
func loadContentList(c *[]content.ContentListItem) error {
	// Open our jsonFile
	jsonFile, err := os.Open("./data/content_all.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		return err
	}

	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	json.Unmarshal(byteValue, &c)

	return nil
}
