package services

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strconv"
	"time"

	"github.com/beowulf-rohan/go-url-shortner/database"
	"github.com/beowulf-rohan/go-url-shortner/model"
	"github.com/go-redis/redis/v8"
)

var (
	config *model.Config
)

func Init(configurations *model.Config) {
	config = configurations
}

func Shorten(request *model.Request, ip string) (*model.Response, error, int) {
	redisClient := database.CreareRedisClient(1)
	defer redisClient.Close()

	_, err := redisClient.Get(database.Ctx, ip).Result()
	if err == redis.Nil {
		_ = redisClient.Set(database.Ctx, ip, config.ApiQuota, 30*time.Minute).Err()
	} else if err != nil {
		val, _ := redisClient.Get(database.Ctx, ip).Result()
		intVal, _ := strconv.Atoi(val)
		if intVal <= 0 {
			limit, _ := redisClient.TTL(database.Ctx, ip).Result()
			return &model.Response{}, fmt.Errorf("rate limit exceeded, rateLimit resets in: %+v", limit/time.Nanosecond/time.Minute), 429
		}
	}

	var shortId string
	if request.ShortenedURL == "" {
		shortId = generateShortCode(request.URL)
	} else {
		shortId = request.ShortenedURL
	}

	redisClient2 := database.CreareRedisClient(0)
	defer redisClient2.Close()

	val, _ := redisClient2.Get(database.Ctx, shortId).Result()
	if val != "" {
		return &model.Response{}, fmt.Errorf("shortened url already in use"), 403
	}

	if request.Expiry == 0 {
		request.Expiry = 30
	}

	err = redisClient2.Set(database.Ctx, shortId, request, request.Expiry*time.Hour).Err()
	if err != nil {
		return &model.Response{}, err, 500
	}

	response := model.Response{
		URL:            request.URL,
		ShortenedURL:   "",
		Expiry:         request.Expiry,
		RateRemaining:  10,
		RateLimitReset: 30,
	}

	redisClient.Decr(database.Ctx, ip)

	val, _ = redisClient2.Get(database.Ctx, ip).Result()
	response.RateRemaining, _ = strconv.Atoi(val)

	ttl, _ := redisClient2.TTL(database.Ctx, ip).Result()
	response.RateLimitReset = ttl / time.Nanosecond / time.Minute

	response.ShortenedURL = config.Domain + "/" + shortId
	return &response, nil, 200
}

func generateShortCode(url string) string {
	hash := sha256.Sum256([]byte(url))
	encoded := base64.URLEncoding.EncodeToString(hash[:])
	return encoded[:8]
}
