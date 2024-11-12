package services

import (
	database "github.com/beowulf-rohan/go-url-shortner/redis"

	"github.com/go-redis/redis/v8"
)

func Resolve(url string) (string, error, int) {
	redisClient := database.CreareRedisClient(0)
	defer redisClient.Close()

	value, err := redisClient.Get(database.Ctx, url).Result()
	if err == redis.Nil {
		return "", err, 404
	} else if err != nil {
		return "", err, 500
	}

	rInr := database.CreareRedisClient(1)
	defer rInr.Close()

	_ = rInr.Incr(database.Ctx, "counter")

	return value, nil, 200
}
