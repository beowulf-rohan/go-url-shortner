package database

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/beowulf-rohan/go-url-shortner/config"
	"github.com/beowulf-rohan/go-url-shortner/model"
	"github.com/go-redis/redis/v8"
)

const (
	URL_METADATA_DB = 0
	IP_DB           = 1
)

type RedisClient struct {
	Client *redis.Client
	DbNo   int
	Ctx    context.Context
}

func GetRedisClient(dbNo int) *RedisClient {
	config := config.GlobalConfig
	reddisClient := redis.NewClient(&redis.Options{
		Addr:     config.DbAddr,
		Password: config.DbPass,
		DB:       dbNo,
	})

	return &RedisClient{
		Client: reddisClient,
		DbNo:   dbNo,
		Ctx:    context.Background(),
	}
}

func (client *RedisClient) PushToRedis(shortenedURLResponse *model.Response) error {
	data, err := json.Marshal(shortenedURLResponse)
	if err != nil {
		log.Printf("Failed to serialize shortenedURLResponse for Redis: %v", err)
		return err
	}

	err = client.Client.Set(client.Ctx, shortenedURLResponse.ShortURL, data, 24*time.Hour).Err()
	if err != nil {
		log.Printf("Failed to push shortenedURLResponse to Redis: %v", err)
		return err
	}

	log.Printf("Successfully pushed shortenedURLResponse %s to Redis with 1-day expiration", shortenedURLResponse.ShortURL)
	return nil
}
