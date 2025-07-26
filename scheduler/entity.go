package scheduler

import (
	"github.com/likecodingloveproblems/sms-gateway/types"
	"time"
)

type Category int

const (
	NormalCategory = iota
	ExpressCategory
)

type Message struct {
	ID        types.ID  `json:"id,omitempty"`
	Text      string    `json:"text,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	Status    string    `json:"status,omitempty"`
}
