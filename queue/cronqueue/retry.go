package cronqueue

import "github.com/herb-go/notification/notificationqueue"

type RetryHandler interface {
	HandleRetry(*notificationqueue.Execution) (bool, error)
}
