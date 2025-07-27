package gateway

import (
	"github.com/likecodingloveproblems/sms-gateway/entity"
	"github.com/likecodingloveproblems/sms-gateway/types"
)

type SendMessageRequest struct {
	UserId    uint64 `json:"user_id"`
	Recipient string `json:"recipient"`
	Text      string `json:"text"`
	Type      string `json:"type"`
}

func mapToMessage(request SendMessageRequest) entity.Message {
	var t entity.MessageType
	switch request.Type {
	case "express":
		t = entity.ExpressMessage
	case "normal":
		t = entity.NormalMessage
	default:
		t = entity.NormalMessage
	}
	return entity.Message{
		UserID:    types.ID(request.UserId),
		Recipient: request.Recipient,
		Text:      request.Text,
		Type:      t,
	}
}
