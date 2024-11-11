package services

import (
	"fmt"
	"log"
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

	value, err := redisClient.Get(database.Ctx, ip).Result()
	if err == redis.Nil {
		_ = redisClient.Set(database.Ctx, ip, config.ApiQuota, 30*time.Minute).Err()
	} else if err != nil {
		val, _ := redisClient.Get(database.Ctx, ip).Result()
		intVal, _ := strconv.Atoi(val)
		if intVal <= 0 {
			limit, _ := redisClient.TTL(database.Ctx, ip).Result()
			return &model.Response{}, fmt.Errorf("rate limit exceeded, rateLimit resets in: %+v", limit), 429
		}
	}
	redisClient.Decr(database.Ctx, ip)
	//TODO: remove this and return ob correctly
	log.Println(value)
	return &model.Response{}, nil, 200
}
