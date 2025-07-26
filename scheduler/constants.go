package scheduler

import "time"

const (
	ProbabilisticSchedulerInterval    = time.Duration(time.Second * 10)
	ExpressStreamKey                  = "sms:express"
	SLAExpressMessageDeliveryDuration = time.Duration(time.Second * 20)
	SLAOptimisticLevel                = 0.7
	SLACriticalLevel                  = 0.9
	SlidingWindowDeliveryDurationSize = 1000
)
