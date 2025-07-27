package entity

import (
	"github.com/likecodingloveproblems/sms-gateway/types"
	"time"
)

type MessageType int

const (
	NormalMessage = iota
	ExpressMessage
)

type MessageStatus string

const (
	QueuedMessage    MessageStatus = "queued"
	DeliveredMessage               = "delivered"
	FailedMessage                  = "failed"
)

type Message struct {
	ID                    types.ID
	UserID                types.ID
	Text                  string
	Recipient             string
	Type                  MessageType
	Status                MessageStatus
	ReceivedAt            time.Time
	DeliveredToOperatorAt time.Time
	FailedAt              time.Time
}
