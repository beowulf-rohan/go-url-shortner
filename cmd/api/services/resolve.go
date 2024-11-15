package services

import (
	"log"

	"github.com/beowulf-rohan/go-url-shortner/config"
	"github.com/beowulf-rohan/go-url-shortner/elasticsearch"
	"github.com/beowulf-rohan/go-url-shortner/model"
	"github.com/beowulf-rohan/go-url-shortner/utils"
	// database "github.com/beowulf-rohan/go-url-shortner/redis"
	// "github.com/go-redis/redis/v8"
)

func Resolve(shortUrl string) (*model.Response, int, error) {
	log.Println("Short URL received for resolution:", shortUrl)

	config := config.GlobalConfig
	ElasticClient, err := elasticsearch.GetElasticClient(config.UrlMetadataIndex)
	if err != nil {
		return &model.Response{}, 500, err
	}

	query := utils.GetResolveQuery(shortUrl)

	existingDoc, err := ElasticClient.GetFromElastic(query)
	if err != nil {
		log.Println("Error checking existing URL:", err)
		return &model.Response{}, 500, err
	}

	log.Printf("shortURL: '%s' mapped to URL: '%s'", shortUrl, existingDoc.URL)
	return existingDoc, 200, nil
}
