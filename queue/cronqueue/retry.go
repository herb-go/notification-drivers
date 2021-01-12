package cronqueue

import (
	"time"

	"github.com/herb-go/notification/notificationqueue"
)

type RetryHandler interface {
	HandleRetry(*notificationqueue.Execution) (bool, error)
}

type PlainRetryHandler []time.Duration

func (h *PlainRetryHandler) HandleRetry(e *notificationqueue.Execution) (bool, error) {
	if e.RetryCount > int32(len(*h)) {
		return false, nil
	}
	now := time.Now()
	e.StartTime = now.Unix()
	e.RetryAfterTime = now.Add((*h)[int(e.RetryCount)]).Unix()
	e.RetryCount = e.RetryCount + 1
	return true, nil
}
