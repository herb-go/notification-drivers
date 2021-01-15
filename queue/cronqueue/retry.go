package cronqueue

import (
	"time"

	"github.com/herb-go/notification/notificationdelivery/notificationqueue"
)

//RetryHandler queue retry handler
type RetryHandler interface {
	//HandleRetry execution update retry count,start time,retry after time.
	//Return false if retry too many.
	HandleRetry(*notificationqueue.Execution) (bool, error)
}

//PlainRetryHandler plain retry handler
type PlainRetryHandler []time.Duration

//HandleRetry execution update retry count,start time,retry after time.
//Return false if retry too many.
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

//PlainRetry plain retry config
type PlainRetry []string

//CreateRetryHandler create retry hanlder
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
