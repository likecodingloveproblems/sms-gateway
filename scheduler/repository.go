package scheduler

import (
	"context"
	"fmt"
	"github.com/likecodingloveproblems/sms-gateway/entity"
	"github.com/likecodingloveproblems/sms-gateway/gateway"
	"github.com/likecodingloveproblems/sms-gateway/pkg/sliding_window"
	"github.com/likecodingloveproblems/sms-gateway/types"
	"github.com/redis/go-redis/v9"
	"log"
	"strconv"
	"strings"
	"time"
)

type RedisRepository struct {
	rdb           *redis.Client
	slidingWindow sliding_window.SlidingWindow
	groupName     string // This must be given to the process by orchestrator
	consumerName  string // This is generated and try to be random and time stamp to create a unique one or given by orchestrator
}

func GetStreamKey(message entity.Message) string {
	switch message.Type {
	case entity.ExpressMessage:
		return ExpressStreamKey
	default:
		// other cases are normal
		return fmt.Sprintf(gateway.NormalStreamKeyTemplate, message.UserID)
	}
}

func (r *RedisRepository) Ack(ctx context.Context, message entity.Message) error {
	return r.rdb.XAck(ctx, GetStreamKey(message), r.groupName, fmt.Sprintf("%d", message.ID)).Err()
}

func NewRepository(rdb *redis.Client) Repository {
	return &RedisRepository{
		rdb:           rdb,
		slidingWindow: sliding_window.NewSlidingWindow(rdb, "sms:express:delivery_sliding_window", SlidingWindowDeliveryDurationSize),
	}
}

func (r *RedisRepository) GetExpressMessagesCount(ctx context.Context) (int64, error) {
	return r.rdb.XLen(ctx, ExpressStreamKey).Result()
}

func (r *RedisRepository) AvgExpressMessageProcessingDuration() (time.Duration, error) {
	avg, err := r.slidingWindow.Avg()
	if err != nil {
		log.Printf("Error in AvgExpressMessageProcessingDuration: %s\n", err.Error())
		return time.Duration(0), err
	}
	return time.Duration(avg), nil
}

func (r *RedisRepository) AddSuccessfulMessageToTimeWindow(message entity.Message) error {
	return r.slidingWindow.Add(float64(time.Now().Second() - message.CreatedAt.Second()))
}

func (r *RedisRepository) Keys(ctx context.Context, pattern string) ([]string, error) {
	// find keys with specified pattern
	return r.rdb.Keys(ctx, pattern).Result()
}

func (r *RedisRepository) ReadStreams(ctx context.Context, streamsKey []string) ([]entity.Message, error) {
	var messages []entity.Message
	streamsArgs := make([]string, 0, len(streamsKey)*2)
	for _, s := range streamsKey {
		streamsArgs = append(streamsArgs, s)
	}
	for range streamsKey {
		streamsArgs = append(streamsArgs, "0-0")
	}

	// List pending
	// XPending -> consumer
	streams, err := r.rdb.XReadGroup(ctx, &redis.XReadGroupArgs{
		Group:    r.groupName,
		Consumer: r.consumerName,
		Streams:  streamsArgs,
		Count:    100,
		//Block:    100 * time.Millisecond,
	}).Result()

	if err != nil && err != redis.Nil {
		log.Printf("[%s] Read error: %v", r.consumerName, err)
		return messages, nil
	}

	for _, stream := range streams {
		for _, message := range stream.Messages {
			messages = append(messages, r.mapRedisMessageToMessage(message))
		}
	}
	return messages, nil
}

func (*RedisRepository) mapRedisMessageToMessage(message redis.XMessage) entity.Message {
	var createdAt time.Time
	if createdAtVal, ok := message.Values["created_at"]; ok {
		if createdAtStr, ok := createdAtVal.(string); ok {
			parsedTime, err := time.Parse(time.RFC3339, createdAtStr)
			if err == nil {
				createdAt = parsedTime
			} else {
				createdAt = time.Now()
			}
		} else {
			createdAt = time.Now()
		}
	} else {
		createdAt = time.Now()
	}

	// Redis Stream IDs are not integers, so parsing to int may not work
	// Use message.ID directly or parse only if you're sure of the format
	// We assume the second part is not used at all and is 0 all the time
	var numericID int
	parts := strings.Split(message.ID, "-")
	if len(parts) > 0 {
		numericID, _ = strconv.Atoi(parts[0]) // best-effort parse
	}

	var text string
	if t, ok := message.Values["text"].(string); ok {
		text = t
	}

	var status string
	if s, ok := message.Values["status"].(string); ok {
		status = s
	}

	return entity.Message{
		ID:        types.ID(numericID),
		Text:      text,
		CreatedAt: createdAt,
		Status:    status,
	}
}

var _ Repository = &RedisRepository{}
