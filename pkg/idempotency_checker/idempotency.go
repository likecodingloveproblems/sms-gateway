package idempotency_checker

import (
	"context"
	"github.com/redis/go-redis/v9"
	"log"
	"time"
)

type IdempotencyChecker interface {
	IsValid(id string) bool
}

type RedisIdempotencyChecker struct {
	rdb        *redis.Client
	ttl        time.Duration
	keyBuilder func(id string) string
}

func NewIdempotencyChecker(rdb *redis.Client, ttl time.Duration, keyBuilder func(string) string) IdempotencyChecker {
	return &RedisIdempotencyChecker{
		rdb:        rdb,
		ttl:        ttl,
		keyBuilder: keyBuilder,
	}
}

func (c RedisIdempotencyChecker) IsValid(id string) bool {
	isSet, err := c.rdb.SetNX(
		context.Background(),
		c.keyBuilder(id),
		true,
		c.ttl,
	).Result()
	if err != nil {
		log.Printf("Error in setnx for %s\n", c.keyBuilder(id))
		return true
	}
	return isSet
}

var _ IdempotencyChecker = RedisIdempotencyChecker{}
