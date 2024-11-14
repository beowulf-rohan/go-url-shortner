package services

import (
	"crypto/sha256"
	"encoding/base64"
	"log"
	"strings"
	"time"

	"github.com/beowulf-rohan/go-url-shortner/config"
	"github.com/beowulf-rohan/go-url-shortner/elasticsearch"
	"github.com/beowulf-rohan/go-url-shortner/model"
)

func Shorten(request *model.Request, ip string) (*model.Response, int, error) {
	log.Printf("Received a request to shorten url: %s", request.URL)
	
	config := config.GlobalConfig
	ElasticClient, err := elasticsearch.GetElasticClient(config.UrlMetadataIndex)
	if err != nil {
		return &model.Response{}, 500, err
	}

	if request.ShortenedURL == "" {
		request.ShortenedURL = generateShortCode(request.URL)
	}
	if request.Expiry == 0 {
		request.Expiry = 24
	}

	existingDoc, err := ElasticClient.GetFromElastic(request.URL)
	if err == nil && existingDoc != nil {
		log.Printf("URL already shortened: %s -> %s", existingDoc.URL, existingDoc.ShortenedURL)
		return existingDoc, 200, nil
	} else if err != nil && !strings.Contains(err.Error(), "not found") {
		log.Println("Error checking existing URL:", err)
		return &model.Response{}, 500, err
	}

	shortenedURLResponse := model.Response{
		URL:          request.URL,
		ShortenedURL: request.ShortenedURL,
		Expiry:       time.Now().Add(request.Expiry * 24 * time.Hour),
	}

	ElasticClient.PushToElastic(shortenedURLResponse, shortenedURLResponse.URL)

	log.Printf("Successfully shortened URL: %s -> %s", shortenedURLResponse.URL, shortenedURLResponse.ShortenedURL)

	return &shortenedURLResponse, 200, nil
}

func generateShortCode(url string) string {
	hash := sha256.Sum256([]byte(url))
	encoded := base64.URLEncoding.EncodeToString(hash[:])
	return encoded[:8]
}
