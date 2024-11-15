package services

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/beowulf-rohan/go-url-shortner/config"
	"github.com/beowulf-rohan/go-url-shortner/elasticsearch"
	"github.com/beowulf-rohan/go-url-shortner/model"
	redisdb "github.com/beowulf-rohan/go-url-shortner/redis"
	"github.com/beowulf-rohan/go-url-shortner/utils"
	"github.com/go-redis/redis/v8"
)

func Resolve(shortUrl string, ip string) (*model.Response, int, error) {
	log.Println("Short URL received for resolution:", shortUrl)

	config := config.GlobalConfig

	isWithinRateLimit, err := CheckIfWithinRateLimit(ip)
	if err != nil {
		log.Println(err)
		return &model.Response{}, 500, err
	}
	if !isWithinRateLimit {
		return &model.Response{}, 429, fmt.Errorf("rate limit exceeded, wait for timer to reset")
	}

	existingDocInRedis, err := CheckIfDocExistInRedis(shortUrl)
	if err != nil {
		log.Println("error checking in redis:", err)
		return &model.Response{}, 500, err
	}

	if existingDocInRedis != nil {
		log.Printf("shortURL: '%s' mapped to URL: '%s', fetched from redis", shortUrl, existingDocInRedis.URL)
		return existingDocInRedis, 200, nil
	}

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

	RedisUrlClient := redisdb.GetRedisClient(redisdb.URL_METADATA_DB)
	err = RedisUrlClient.PushToRedis(existingDoc)
	if err != nil {
		log.Println("error publishing to redis", err)
		return &model.Response{}, 500, err
	}

	log.Printf("shortURL: '%s' mapped to URL: '%s'", shortUrl, existingDoc.URL)
	return existingDoc, 200, nil
}

func CheckIfWithinRateLimit(ip string) (bool, error) {
	RedisIpClient := redisdb.GetRedisClient(redisdb.IP_DB)
	config := config.GlobalConfig

	_, err := RedisIpClient.Client.Get(RedisIpClient.Ctx, ip).Result()
	if err == redis.Nil {
		_ = RedisIpClient.Client.Set(RedisIpClient.Ctx, ip, config.ApiRateLimit, time.Hour)
	} else if err != nil {
		return false, err
	}

	_ = RedisIpClient.Client.Decr(RedisIpClient.Ctx, ip)
	val, _ := RedisIpClient.Client.Get(RedisIpClient.Ctx, ip).Result()
	ttl := RedisIpClient.Client.TTL(RedisIpClient.Ctx, ip).Val().Minutes()
	
	valInt, _ := strconv.Atoi(val)
	if valInt <= 0 {
		return false, nil
	}

	fmt.Printf("%s: requests left, time before counter resets: %.2f minutes", val, ttl)

	return true, nil
}

func CheckIfDocExistInRedis(shortUrl string) (*model.Response, error) {
	RedisUrlClient := redisdb.GetRedisClient(redisdb.URL_METADATA_DB)

	data, err := RedisUrlClient.Client.Get(RedisUrlClient.Ctx, shortUrl).Result()
	if err == redis.Nil {
		log.Printf("No data found for short URL in redis: %s", shortUrl)
		return nil, nil
	} else if err != nil {
		log.Printf("Error fetching data for short URL %s from Redis: %v", shortUrl, err)
		return nil, err
	}

	var response model.Response
	err = json.Unmarshal([]byte(data), &response)
	if err != nil {
		log.Printf("Failed to unmarshal data for short URL %s: %v", shortUrl, err)
		return nil, err
	}

	return &response, nil
}
