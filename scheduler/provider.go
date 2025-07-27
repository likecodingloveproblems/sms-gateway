package scheduler

import (
	"context"
	"errors"
	"github.com/likecodingloveproblems/sms-gateway/entity"
	"log"
	"math/rand"
	"time"
)

type Provider interface {
	Send(ctx context.Context, message entity.Message) error
}

type LogProvider struct{}

func (l LogProvider) Send(ctx context.Context, message entity.Message) error {
	log.Printf("sent message: %v\n", message)
	return nil
}

var _ Provider = LogProvider{}

type RandomlyFailProviderWithDelay struct {
	FailurePerc int // between 0 and 100
	MaxDelay    time.Duration
}

func (r RandomlyFailProviderWithDelay) Send(ctx context.Context, message entity.Message) error {
	delay := time.Duration(rand.Intn(int(r.MaxDelay.Milliseconds())))
	time.Sleep(delay)
	if rand.Intn(100)+1 < r.FailurePerc {
		return errors.New("Failed")
	}
	return nil
}

var _ Provider = RandomlyFailProviderWithDelay{}
