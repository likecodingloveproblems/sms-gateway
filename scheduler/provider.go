package scheduler

import (
	"errors"
	"log"
	"math/rand"
	"time"
)

type Provider interface {
	Send(message Message) error
}

type LogProvider struct{}

func (l LogProvider) Send(message Message) error {
	log.Printf("sent message: %v\n", message)
	return nil
}

var _ Provider = LogProvider{}

type RandomlyFailProviderWithDelay struct {
	FailurePerc int // between 0 and 100
	MaxDelay    time.Duration
}

func (r RandomlyFailProviderWithDelay) Send(message Message) error {
	delay := time.Duration(rand.Intn(int(r.MaxDelay.Milliseconds())))
	time.Sleep(delay)
	if rand.Intn(100)+1 < r.FailurePerc {
		return errors.New("Failed")
	}
	return nil
}

var _ Provider = RandomlyFailProviderWithDelay{}
