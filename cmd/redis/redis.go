package database

import (
	"context"

	"github.com/beowulf-rohan/go-url-shortner/config"
	"github.com/go-redis/redis/v8"
)

var Ctx = context.Background()

func CreareRedisClient(dbNo int) *redis.Client {
	config := config.GlobalConfig
	reddisClient := redis.NewClient(&redis.Options{
		Addr:     config.DbAddr,
		Password: config.DbPass,
		DB:       dbNo,
	})
	return reddisClient
}
