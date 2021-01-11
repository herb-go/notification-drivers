package cronqueue

import (
	"github.com/herb-go/notification/notificationqueue"
)

type Store interface {
	List(start string, count int) ([]*notificationqueue.Execution, error)
	Insert(execution *notificationqueue.Execution) error
	Replace(eid string, new *notificationqueue.Execution) error
	Remove(nid string) error
	Start() error
	Stop() error
}
