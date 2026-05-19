package goredis

import (
	"time"

	"github.com/go-redis/redis/v8"
)

var rclient *redis.Client

func init() {
	InitRedisClient()
}
func InitRedisClient() {
	rclient = redis.NewClient(&redis.Options{
		Addr:        "localhost:6379",
		Password:    "",
		DB:          0,
		DialTimeout: 1 * time.Second,
	})
}
func GetRedisClient() *redis.Client {
	return rclient
}
