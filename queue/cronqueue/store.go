package cronqueue

import (
	"github.com/herb-go/notification/notificationdelivery/notificationqueue"
)

//Store queue store interface
type Store interface {
	//List list queued Execution form start
	List(start string, count int) ([]*notificationqueue.Execution, string, error)
	//Insert insert execution to queue.
	//Do nothing if notifiaction exsits
	Insert(execution *notificationqueue.Execution) error
	//Replace exectution in queue with given eid.
	//Do nothing if eid not match.
	Replace(eid string, new *notificationqueue.Execution) error
	//Remove notifcation with given nid form queue.
	//Do nothing if notificatio not found
	Remove(nid string) error
	//Start start store
	Start() error
	//Stop stop store
	Stop() error
}
