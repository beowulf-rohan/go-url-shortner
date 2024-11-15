package services

import (
	"log"

	database "github.com/beowulf-rohan/go-url-shortner/redis"

	"github.com/go-redis/redis/v8"
)

func Resolve(url string) (string, int, error) {
	log.Println("Short URL received for resolution:", url)

	redisClient := database.CreareRedisClient(0)
	defer redisClient.Close()

	value, err := redisClient.Get(database.Ctx, url).Result()
	if err == redis.Nil {
		return "", 404, err
	} else if err != nil {
		return "", 500, err
	}

	rInr := database.CreareRedisClient(1)
	defer rInr.Close()

	_ = rInr.Incr(database.Ctx, "counter")

	return value, 200, nil
}
