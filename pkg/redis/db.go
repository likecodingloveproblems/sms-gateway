package redis

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"os"
)

func NewClient(db int) *redis.Client {
	// It was better to use koanf for loading configs
	Host := os.Getenv("REDIS_HOST")
	Port := os.Getenv("REDIS_PORT")
	User := os.Getenv("REDIS_USER")
	Password := os.Getenv("REDIS_PASSWORD")
	return redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", Host, Port),
		Username: User,
		Password: Password,
		DB:       db,
	})
}
