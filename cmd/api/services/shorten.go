package services

import (
	"log"
	"strings"
	"time"

	"github.com/beowulf-rohan/go-url-shortner/config"
	"github.com/beowulf-rohan/go-url-shortner/elasticsearch"
	"github.com/beowulf-rohan/go-url-shortner/model"
	"github.com/beowulf-rohan/go-url-shortner/utils"
)

func Shorten(request *model.Request, ip string) (*model.Response, int, error) {
	log.Printf("Received a request to shorten url: %s", request.URL)

	config := config.GlobalConfig
	ElasticClient, err := elasticsearch.GetElasticClient(config.UrlMetadataIndex)
	if err != nil {
		return &model.Response{}, 500, err
	}

	if request.ShortURL == "" {
		request.ShortURL = utils.GenerateShortUrl(request.URL)
	}
	if request.Expiry == 0 {
		request.Expiry = 24
	}

	query := utils.GetShortenQuery(request.URL)

	existingDoc, err := ElasticClient.GetFromElastic(query)
	if err == nil && existingDoc != nil {
		log.Printf("URL already shortened: %s -> %s", existingDoc.URL, existingDoc.ShortURL)
		return existingDoc, 200, nil
	} else if err != nil && !strings.Contains(err.Error(), "not found") {
		log.Println("Error checking existing URL:", err)
		return &model.Response{}, 500, err
	}

	shortenedURLResponse := model.Response{
		URL:      request.URL,
		ShortURL: request.ShortURL,
		CreatedAt: time.Now(),
		Expiry:   time.Now().Add(request.Expiry * 24 * time.Hour),
	}

	ElasticClient.PushToElastic(shortenedURLResponse, shortenedURLResponse.URL)

	log.Printf("Successfully shortened URL: %s -> %s", shortenedURLResponse.URL, shortenedURLResponse.ShortURL)

	return &shortenedURLResponse, 200, nil
}
