package db

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"os"
)

var ctx = context.Background()

func ConnectRedis() *redis.Client {
	end := godotenv.Load()
	if end != nil {
		panic("Failed to load .env file")
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("RedisHost") + ":" + os.Getenv("RedisPort"),
		Password: os.Getenv("RedisPass"),
		DB:       1,
	})
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		panic("Could not connect to Redis !")
	}

	return rdb
}
