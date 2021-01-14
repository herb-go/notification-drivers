package cronqueue

import (
	"time"

	"github.com/herb-go/notification/notificationdelivery/notificationqueue"
)

type RetryHandler interface {
	HandleRetry(*notificationqueue.Execution) (bool, error)
}

type PlainRetryHandler []time.Duration

func (h *PlainRetryHandler) HandleRetry(e *notificationqueue.Execution) (bool, error) {
	if e.RetryCount >= int32(len(*h)) {
		return false, nil
	}
	now := time.Now()
	e.StartTime = now.Unix()
	e.RetryAfterTime = now.Add((*h)[int(e.RetryCount)]).Unix()
	e.RetryCount = e.RetryCount + 1
	return true, nil
}

type PlainRetry []string

func (r *PlainRetry) CreateRetryHandler() (*PlainRetryHandler, error) {
	var err error
	p := make(PlainRetryHandler, len(*r))
	for k, v := range *r {
		p[k], err = time.ParseDuration(v)
		if err != nil {
			return nil, err
		}
	}
	return &p, err
}
