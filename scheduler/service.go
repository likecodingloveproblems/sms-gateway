package scheduler

import (
	"context"
	"github.com/likecodingloveproblems/sms-gateway/entity"
	"log"
	"math/rand"
	"time"
)

type Repository interface {
	GetExpressMessagesCount(context.Context) (int64, error)
	AvgExpressMessageProcessingDuration() (time.Duration, error)
	AddSuccessfulMessageToTimeWindow(message entity.Message) error
	Keys(ctx context.Context, pattern string) ([]string, error)
	ReadStreams(ctx context.Context, keys []string) ([]entity.Message, error)
	Ack(ctx context.Context, message entity.Message) error
}

type Scheduler interface {
	Run(ctx context.Context) // entry point
	OnSuccess(ctx context.Context, message entity.Message)
	OnFailure(ctx context.Context, message entity.Message)
}

type Worker interface {
	Submit(task func())
	Stop()
}

type ProbabilisticProportionalScheduler struct {
	repository                             Repository
	estimateExpressMessageDeliveryDuration time.Duration
	worker                                 Worker
	provider                               Provider
}

func (s *ProbabilisticProportionalScheduler) OnSuccess(ctx context.Context, message entity.Message) {
	// Now you must ACK
	// it seems that removing message from streams is not responsiblity of service
	// as it's possible to replace redis.streams with kafka.topics
	// but here it seems we want

	// And we must raise and event to be used with other services like reporting and accounting
}

func (s *ProbabilisticProportionalScheduler) OnFailure(ctx context.Context, message entity.Message) {
	// Retry first
	// Then ACK the message and send the message to dead letter queue
	// Then add failure event
	err := s.repository.Ack(ctx, message)
	if err != nil {
		// It's very bad things someone must come and review the process
		log.Printf("Error in Ack: %s", err.Error())
	}

	// Now send the new status of the message to be processed by reporting
	s.repository.Add(ctx, message)
}

func (s *ProbabilisticProportionalScheduler) Run(ctx context.Context) {
	// It's blocking!
	// Schedule to listen on what kind of message
	var streamsKey []string
	var err error
	defer s.worker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		switch s.schedule() {
		case entity.NormalMessage:
			streamsKey, err = s.getNormalStreamsKeyToConsume(ctx)
			if err != nil {
				log.Printf("Error in getStreamsToConsume: %s\n", err.Error())
				continue
			}
		case entity.ExpressMessage:
			streamsKey = []string{ExpressStreamKey}
		}
		messages, err := s.repository.ReadStreams(ctx, streamsKey)
		if err != nil {
			log.Printf("Error in ReadStream: %s\n", err.Error())
			continue
		}
		s.sendToProvider(ctx, messages)
	}
}

func (s *ProbabilisticProportionalScheduler) schedule() entity.MessageType {
	// Now we want to implement a probabilistic proportional sharing algorithm
	// that calculate processing time of all express messages
	// then try to optimize processing time to be 70% expected express processing time that is SLA
	expressPortion := 50
	if s.estimateExpressMessageDeliveryDuration == 0 {
		return entity.NormalMessage
	}
	expressSLARate := s.estimateExpressMessageDeliveryDuration.Seconds() / SLAExpressMessageDeliveryDuration.Seconds()
	if expressSLARate > 0.9 {
		// Status of express is critical
		// keep 5% for normal to prevent starvation
		expressPortion = 95
	} else if expressSLARate > 0.85 {
		expressPortion = 90
	} else if expressSLARate > 0.8 {
		expressPortion = 80
	} else if expressSLARate > 0.75 {
		expressPortion = 70
	} else if expressSLARate > 0.65 {
		// Optimize for 70% of SLA
		expressPortion = 50
	} else if expressSLARate > 0.5 {
		expressPortion = 40
	} else if expressSLARate > 0.3 {
		expressPortion = 30
	} else {
		expressPortion = 20
	}

	if rand.Intn(100) <= expressPortion {
		return entity.ExpressMessage
	}
	return entity.NormalMessage
}

func (s *ProbabilisticProportionalScheduler) updateState(ctx context.Context) {
	// update stream info from redis periodic
	ticker := time.NewTicker(ProbabilisticSchedulerInterval)
	for {
		<-ticker.C
		log.Println("Going to update probabilistic info from redis!")
		expressMessagesCount, err := s.repository.GetExpressMessagesCount(ctx)
		if err != nil {
			log.Printf("Error in XLen express stream: %s\n", err.Error())
			// keep things smooth
			s.estimateExpressMessageDeliveryDuration = time.Duration(SLAExpressMessageDeliveryDuration.Seconds() * 0.7)
			return
		}
		avgExpressMessageProcessingDuration, err := s.repository.AvgExpressMessageProcessingDuration()
		if err != nil {
			log.Printf("Error in AvgExpressMessagesProcessingDuration: %s\n", err.Error())
			// keep things smooth
			s.estimateExpressMessageDeliveryDuration = time.Duration(SLAExpressMessageDeliveryDuration.Seconds() * 0.7)
			return
		}
		s.estimateExpressMessageDeliveryDuration = time.Duration(expressMessagesCount) * avgExpressMessageProcessingDuration
	}
}

func (s *ProbabilisticProportionalScheduler) getNormalStreamsKeyToConsume(ctx context.Context) ([]string, error) {
	// It can be part of the broker, to load balance between partitions of normal sms

	// What about a time that number of streams is huge
	// In the problem specification is said we will have only 100,000 customers that is not a huge number
	// but if the number of customers go up by order we can implement a load balancer here for streams
	// like round-robin
	keys, err := s.repository.Keys(ctx, "sms:normal:*")
	if err != nil {
		return []string{}, err
	}
	return keys, nil
}

func (s *ProbabilisticProportionalScheduler) sendToProvider(ctx context.Context, messages []entity.Message) {
	// It must send to provider with onSuccess and onFailure callbacks
	// All of this will be submitted to a worker to keep the backpressure on control
	// a worker is used with limited amount of concurrency
	for _, message := range messages {
		// It's better to have a retry policy hear on failure
		if err := s.provider.Send(ctx, message); err != nil {
			s.OnFailure(ctx, message)
		} else {
			s.OnSuccess(ctx, message)
		}
	}
}

func NewScheduler(repository Repository, worker Worker, provider Provider) Scheduler {
	return &ProbabilisticProportionalScheduler{
		repository: repository,
		worker:     worker,
		provider:   provider,
	}
}

var _ Scheduler = &ProbabilisticProportionalScheduler{}
