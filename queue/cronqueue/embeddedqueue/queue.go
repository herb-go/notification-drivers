package embeddedqueue

import (
	"sync"

	"github.com/vmihailenco/msgpack"

	"github.com/herb-go/herbdata"

	"github.com/herb-go/herbdata/kvdb"
	"github.com/herb-go/notification/notificationdelivery/notificationqueue"
)

//Engine embedded store
type Engine struct {
	locker sync.Mutex
	DB     *kvdb.Database
}

//List list queued Execution form start
func (e *Engine) List(start string, count int) ([]*notificationqueue.Execution, string, error) {
	e.locker.Lock()
	defer e.locker.Unlock()
	iter := []byte(start)
	data, iter, err := e.DB.Next(iter, count)
	if err != nil {
		return nil, "", nil
	}
	var results = make([]*notificationqueue.Execution, len(data))
	for k := range data {
		results[k] = notificationqueue.NewExecution()
		err = msgpack.Unmarshal(data[k].Value, results[k])
		if err != nil {
			return nil, "", nil
		}
	}
	return results, string(iter), nil
}

//Insert insert execution to queue.
//Do nothing if notifiaction exsits
func (e *Engine) Insert(execution *notificationqueue.Execution) error {
	e.locker.Lock()
	defer e.locker.Unlock()
	bs, err := msgpack.Marshal(execution)
	if err != nil {
		return err
	}
	_, err = e.DB.Get([]byte(execution.Notification.ID))
	if err == nil {
		return nil
	}
	if err != herbdata.ErrNotFound {
		return err
	}
	return e.DB.Set([]byte(execution.Notification.ID), bs)
}

//Replace exectution in queue with given eid.
//Do nothing if eid not match.
func (e *Engine) Replace(eid string, new *notificationqueue.Execution) error {
	e.locker.Lock()
	defer e.locker.Unlock()
	bs, err := e.DB.Get([]byte(new.Notification.ID))
	if err != nil {
		if err == herbdata.ErrNotFound {
			return nil
		}
		return err
	}
	execution := notificationqueue.NewExecution()
	err = msgpack.Unmarshal(bs, execution)
	if err != nil {
		return err
	}
	if execution.ExecutionID != eid {
		return nil
	}
	bs, err = msgpack.Marshal(new)
	if err != nil {
		return err
	}
	return e.DB.Set([]byte(new.Notification.ID), bs)
}

//Remove notifcation with given nid form queue.
//Do nothing if notificatio not found
func (e *Engine) Remove(nid string) error {
	e.locker.Lock()
	defer e.locker.Unlock()
	return e.DB.Delete([]byte(nid))
}

//Start start store
func (e *Engine) Start() error {
	return e.DB.Start()
}

//Stop stop store
func (e *Engine) Stop() error {
	return e.DB.Stop()
}

//New create embedded store
func New() *Engine {
	return &Engine{}
}
