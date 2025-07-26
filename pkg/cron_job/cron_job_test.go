package cron_job_test

import (
	"context"
	"github.com/likecodingloveproblems/sms-gateway/pkg/cron_job"
	"sync"
	"testing"
	"time"
)

func TestCronJob_Run(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var mu sync.Mutex
	count := 0
	triggered := make(chan struct{}, 5)

	jobFunc := func(ctx context.Context) {
		mu.Lock()
		count++
		mu.Unlock()
		triggered <- struct{}{}
	}

	cron := cron_job.NewCronJob(ctx, 10*time.Millisecond, jobFunc)

	// Run in background
	go cron.Run()

	// Wait for a few job executions
	timeout := time.After(100 * time.Millisecond)
LOOP:
	for {
		select {
		case <-triggered:
			if count >= 3 {
				break LOOP
			}
		case <-timeout:
			t.Fatal("timeout waiting for job to trigger")
		}
	}

	// Cancel the context and wait a bit
	cancel()
	time.Sleep(20 * time.Millisecond)

	// Capture current count
	mu.Lock()
	finalCount := count
	mu.Unlock()

	// Wait more to ensure no more jobs run
	time.Sleep(30 * time.Millisecond)

	mu.Lock()
	if count > finalCount {
		t.Errorf("expected job to stop after context cancel, got extra calls: before=%d after=%d", finalCount, count)
	}
	mu.Unlock()
}
