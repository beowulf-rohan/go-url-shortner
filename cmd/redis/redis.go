package database

import (
	"context"

	"github.com/beowulf-rohan/go-url-shortner/model"
	"github.com/go-redis/redis/v8"
)

var (
	config *model.Config
)

var Ctx = context.Background()

func Init(configurations *model.Config) {
	config = configurations
}

func CreareRedisClient(dbNo int) *redis.Client {
	reddisClient := redis.NewClient(&redis.Options{
		Addr:     config.DbAddr,
		Password: config.DbPass,
		DB:       dbNo,
	})
	return reddisClient
}
