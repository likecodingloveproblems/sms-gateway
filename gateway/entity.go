package gateway

import "github.com/likecodingloveproblems/sms-gateway/types"

type MessageType int

const (
	NormalMessage = iota
	ExpressMessage
)

type Message struct {
	ID       types.ID
	Text     string
	Receiver string
	Type     MessageType
}
