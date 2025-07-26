package idempotency_checker_test

import (
	"github.com/likecodingloveproblems/sms-gateway/pkg/idempotency_checker"
	"github.com/likecodingloveproblems/sms-gateway/tests"
	"testing"
	"time"
)

func TestRedisIdempotencyChecker_isValid(t *testing.T) {
	rdb := tests.NewTestRedisClient()

	checker := idempotency_checker.NewIdempotencyChecker(rdb, time.Second, func(id string) string {
		return "idem:" + id
	})

	id := "12345"

	// First time should return true (new id)
	if !checker.IsValid(id) {
		t.Error("expected true on first call")
	}

	// Second time should return false (already exists)
	if checker.IsValid(id) {
		t.Error("expected false on second call")
	}

	// Wait for ttl to expire
	time.Sleep(2 * time.Second)

	// After expiration, should return true again
	if !checker.IsValid(id) {
		t.Error("expected true after ttl expired")
	}
}
