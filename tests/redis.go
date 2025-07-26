package tests

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"os"
)

func NewTestRedisClient() *redis.Client {
	addr := os.Getenv("REDIS_TEST_ADDR")
	if addr == "" {
		addr = "localhost:6379"
	}

	db := 15
	if dbEnv := os.Getenv("REDIS_TEST_DB"); dbEnv != "" {
		fmt.Sscanf(dbEnv, "%d", &db)
	}

	return redis.NewClient(&redis.Options{
		Addr: addr,
		DB:   db,
	})
}
