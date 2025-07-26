package cron_job

import (
	"context"
	"log"
	"time"
)

type CronJob struct {
	ctx      context.Context
	Interval time.Duration
	JobFunc  func(ctx context.Context)
}

func NewCronJob(ctx context.Context, interval time.Duration, jobFunc func(ctx2 context.Context)) *CronJob {
	return &CronJob{
		ctx:      ctx,
		Interval: interval,
		JobFunc:  jobFunc,
	}
}

func (c CronJob) Run() {
	ticker := time.NewTicker(c.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-c.ctx.Done():
			log.Println("Cache sync stopped due to context cancellation")
			return
		case <-ticker.C:
			go c.JobFunc(c.ctx)
		}
	}

}
