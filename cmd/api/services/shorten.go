package services

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	redisDB "github.com/beowulf-rohan/go-url-shortner/redis"

	"github.com/beowulf-rohan/go-url-shortner/model"
	"github.com/go-redis/redis/v8"
)

var (
	config *model.Config
)

func Init(configurations *model.Config) {
	config = configurations
}

func Shorten(request model.Request, ip string) (*model.Response, error, int) {
	redisClient := redisDB.CreareRedisClient(1)
	defer redisClient.Close()

	val, err := redisClient.Get(redisDB.Ctx, ip).Result()
	if err == redis.Nil {
		err = redisClient.Set(redisDB.Ctx, ip, config.ApiQuota, 30*60*time.Second).Err()
		if err != nil {
			return &model.Response{}, err, 500
		}
	} else {
		intVal, _ := strconv.Atoi(val)
		if intVal <= 0 {
			limit, _ := redisClient.TTL(redisDB.Ctx, ip).Result()
			return &model.Response{}, fmt.Errorf("rate limit exceeded, rateLimit resets in: %+v", limit/time.Nanosecond/time.Minute), 429
		}
	}

	var shortId string
	if request.ShortenedURL == "" {
		shortId = generateShortCode(request.URL)
	} else {
		shortId = request.ShortenedURL
	}

	redisClient2 := redisDB.CreareRedisClient(0)
	defer redisClient2.Close()

	val, _ = redisClient2.Get(redisDB.Ctx, shortId).Result()
	if val != "" {
		return &model.Response{}, fmt.Errorf("shortened url already in use"), 403
	}

	if request.Expiry == 0 {
		request.Expiry = 30
	}

	log.Printf("error check here....%+v", request)
	jsonData, err := json.Marshal(request.ShortenedURL)
	if err != nil {
		return &model.Response{}, fmt.Errorf("failed to serialize request: %v", err), 500
	}
	err = redisClient2.Set(redisDB.Ctx, shortId, jsonData, request.Expiry*time.Hour).Err()
	if err != nil {
		return &model.Response{}, err, 500
	}
	log.Println("error check pass.....")

	response := model.Response{
		URL:            request.URL,
		ShortenedURL:   "",
		Expiry:         request.Expiry,
		RateRemaining:  10,
		RateLimitReset: 30,
	}

	redisClient.Decr(redisDB.Ctx, ip)

	val, _ = redisClient2.Get(redisDB.Ctx, ip).Result()
	response.RateRemaining, _ = strconv.Atoi(val)

	ttl, _ := redisClient2.TTL(redisDB.Ctx, ip).Result()
	response.RateLimitReset = ttl / time.Nanosecond / time.Minute

	response.ShortenedURL = config.Domain + "/" + shortId
	return &response, nil, 200
}

func generateShortCode(url string) string {
	hash := sha256.Sum256([]byte(url))
	encoded := base64.URLEncoding.EncodeToString(hash[:])
	return encoded[:8]
}
