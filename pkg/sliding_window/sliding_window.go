package sliding_window

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"strconv"
)

type SlidingWindow interface {
	Add(item float64) error
	Avg() (float64, error)
}

type RedisListSlidingWindow struct {
	rdb  *redis.Client
	key  string
	size int64
}

func (r RedisListSlidingWindow) Add(item float64) error {
	// items can be stashed and sent once to redis
	return r.rdb.LPush(context.Background(), r.key, item).Err()
}

func (r RedisListSlidingWindow) Avg() (float64, error) {
	ctx := context.Background()
	pipe := r.rdb.Pipeline()
	pipe.LTrim(ctx, r.key, 0, r.size)
	valuesCmd := pipe.LRange(ctx, r.key, 0, r.size)
	_, err := pipe.Exec(ctx)
	if err != nil {
		return 0, err
	}
	values, err := valuesCmd.Result()
	if err != nil {
		return 0, fmt.Errorf("lrange failed: %w", err)
	}

	if len(values) == 0 {
		return 0, nil // no values to average
	}

	var sum float64
	for _, v := range values {
		num, err := strconv.ParseFloat(v, 64)
		if err != nil {
			continue // skip invalid entries
		}
		sum += num
	}

	avg := sum / float64(len(values))
	return avg, nil
}

func NewSlidingWindow(rdb *redis.Client, key string, size int64) SlidingWindow {
	return &RedisListSlidingWindow{
		rdb:  rdb,
		key:  key,
		size: size,
	}
}

var _ SlidingWindow = &RedisListSlidingWindow{}
