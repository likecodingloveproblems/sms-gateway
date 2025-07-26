package scheduler

import (
	"context"
	"log"
	"math/rand"
	"time"
)

type Repository interface {
	GetExpressMessagesCount(context.Context) (int64, error)
	AvgExpressMessageProcessingDuration() (time.Duration, error)
	AddSuccessfulMessageToTimeWindow(message Message) error
	Keys(ctx context.Context, pattern string) ([]string, error)
	ReadStreams(ctx context.Context, keys []string) ([]Message, error)
}

type Scheduler interface {
	Run(ctx context.Context) // entry point
	OnSuccess(ctx context.Context, message Message)
	OnFailure(ctx context.Context, message Message)
}

type Task interface {
	OnSuccess(ctx context.Context)
	OnFailure(ctx context.Context)
	Run(ctx context.Context)
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

func (p *ProbabilisticProportionalScheduler) OnSuccess(ctx context.Context, message Message) {
	// Now you must ACK
	// it seems that removing message from streams is not responsiblity of service
	// as it's possible to replace redis.streams with kafka.topics
	// but here it seems we want

	// And we must raise and event to be used with other services like reporting and accounting
}

func (p *ProbabilisticProportionalScheduler) OnFailure(ctx context.Context, message Message) {
	// Retry first
	// Then ACK the message and send the message to dead letter queue
	// Then Emit failure event
}

func (p *ProbabilisticProportionalScheduler) Run(ctx context.Context) {
	// It's blocking!
	// Schedule to listen on what kind of message
	var streamsKey []string
	var err error
	defer p.worker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		switch p.schedule() {
		case NormalCategory:
			streamsKey, err = p.getNormalStreamsKeyToConsume(ctx)
			if err != nil {
				log.Printf("Error in getStreamsToConsume: %s\n", err.Error())
				continue
			}
		case ExpressCategory:
			streamsKey = []string{ExpressStreamKey}
		}
		messages, err := p.repository.ReadStreams(ctx, streamsKey)
		if err != nil {
			log.Printf("Error in ReadStream: %s\n", err.Error())
			continue
		}
		p.sendToProvider(ctx, messages)
	}
}

func (p *ProbabilisticProportionalScheduler) schedule() Category {
	// Now we want to implement a probabilistic proportional sharing algorithm
	// that calculate processing time of all express messages
	// then try to optimize processing time to be 70% expected express processing time that is SLA
	expressPortion := 50
	if p.estimateExpressMessageDeliveryDuration == 0 {
		return NormalCategory
	}
	expressSLARate := p.estimateExpressMessageDeliveryDuration.Seconds() / SLAExpressMessageDeliveryDuration.Seconds()
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
		return ExpressCategory
	}
	return NormalCategory
}

func (p *ProbabilisticProportionalScheduler) updateState(ctx context.Context) {
	// update stream info from redis periodic
	ticker := time.NewTicker(ProbabilisticSchedulerInterval)
	for {
		<-ticker.C
		log.Println("Going to update probabilistic info from redis!")
		expressMessagesCount, err := p.repository.GetExpressMessagesCount(ctx)
		if err != nil {
			log.Printf("Error in XLen express stream: %s\n", err.Error())
			// keep things smooth
			p.estimateExpressMessageDeliveryDuration = time.Duration(SLAExpressMessageDeliveryDuration.Seconds() * 0.7)
			return
		}
		avgExpressMessageProcessingDuration, err := p.repository.AvgExpressMessageProcessingDuration()
		if err != nil {
			log.Printf("Error in AvgExpressMessagesProcessingDuration: %s\n", err.Error())
			// keep things smooth
			p.estimateExpressMessageDeliveryDuration = time.Duration(SLAExpressMessageDeliveryDuration.Seconds() * 0.7)
			return
		}
		p.estimateExpressMessageDeliveryDuration = time.Duration(expressMessagesCount) * avgExpressMessageProcessingDuration
	}
}

func (p *ProbabilisticProportionalScheduler) getNormalStreamsKeyToConsume(ctx context.Context) ([]string, error) {
	// It can be part of the broker, to load balance between partitions of normal sms

	// What about a time that number of streams is huge
	// In the problem specification is said we will have only 100,000 customers that is not a huge number
	// but if the number of customers go up by order we can implement a load balancer here for streams
	// like round-robin
	keys, err := p.repository.Keys(ctx, "sms:normal:*")
	if err != nil {
		return []string{}, err
	}
	return keys, nil
}

func (p *ProbabilisticProportionalScheduler) sendToProvider(ctx context.Context, messages []Message) {
	// It must send to provider with onSuccess and onFailure callbacks
	// All of this will be submitted to a worker to keep the backpressure on control
	// a worker is used with limited amount of concurrency

}

func NewScheduler(repository Repository, worker Worker, provider Provider) Scheduler {
	return &ProbabilisticProportionalScheduler{
		repository: repository,
		worker:     worker,
		provider:   provider,
	}
}

var _ Scheduler = &ProbabilisticProportionalScheduler{}
