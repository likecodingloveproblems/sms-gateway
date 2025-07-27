package gateway

import (
	"github.com/likecodingloveproblems/sms-gateway/entity"
	"github.com/redis/go-redis/v9"
)

type Repository struct {
	rdb *redis.Client
}

func NewRepository(rdb *redis.Client) *Repository {
	return &Repository{
		rdb: rdb,
	}
}

func (r Repository) GetMessageUnitPrice() uint64 {
	return DefaultMessagePerUnitPrice
}

func (r Repository) AddMessage(streamKey string, message entity.Message) error {
	// it will add message to the redis stream base on the message type
	return nil
}
