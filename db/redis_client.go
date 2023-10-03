package db

import (
	"github.com/go-redis/redis/v8"
)

func RedisClient() *redis.Client {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "redis-17035.c295.ap-southeast-1-1.ec2.cloud.redislabs.com:17035",
		Password: "qBK3RJEnYyqEklkRCqsyzQ2Q0fvw97eF",
	})

	return redisClient
}
