package db

import "github.com/redis/go-redis/v9"

var (
	Rds *redis.Client
)

func NewRedisClient()  {
	Rds = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password
		DB:       0,  // use default DB
		Protocol: 2,
	})
}
