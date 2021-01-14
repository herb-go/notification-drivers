package cronqueue

import (
	"time"

	"github.com/herb-go/notification"
	"github.com/herb-go/notification/notificationdelivery/notificationqueue"
)

var DefaultInterval = time.Minute
var DefaultExecuteCount = 10

type Queue struct {
	Store            Store
	Timeout          time.Duration
	Interval         time.Duration
	ExecuteCount     int
	IDGenerator      func() (string, error)
	c                chan int
	pipe             chan *notificationqueue.Execution
	OnError          func(error)
	OnDeliverTimeout func(*notificationqueue.Execution)
	OnRetryTooMany   func(*notificationqueue.Execution)
	RetryHandler     RetryHandler
}

func (q *Queue) Recover() {
	r := recover()
	if r != nil {
		err := r.(error)
		q.OnError(err)
	}

}
func (q *Queue) PopChan() (<-chan *notificationqueue.Execution, error) {
	return q.pipe, nil
}

func (q *Queue) NewExecution(n *notification.Notification) (*notificationqueue.Execution, error) {
	var err error
	eid, err := q.IDGenerator()
	if err != nil {
		return nil, err
	}
	e := notificationqueue.NewExecution()
	e.Notification = n
	e.ExecutionID = eid
	e.RetryCount = 0
	return e, nil
}

func (q *Queue) pushExecution(e *notificationqueue.Execution) {
	select {
	case q.pipe <- e:
	case <-time.After(q.Timeout):
		go q.OnDeliverTimeout(e)
	}
}
func (q *Queue) Push(n *notification.Notification) error {
	e, err := q.NewExecution(n)
	if err != nil {
		return err
	}
	err = q.Store.Insert(e)
	if err != nil {
		return err
	}
	go q.pushExecution(e)
	return nil
}
func (q *Queue) Remove(nid string) error {
	return q.Store.Remove(nid)
}

func (q *Queue) retry(e *notificationqueue.Execution) {
	defer q.Recover()
	if time.Now().Unix() <= e.RetryAfterTime {
		return
	}
	eid := e.ExecutionID
	ok, err := q.RetryHandler.HandleRetry(e)
	if err != nil {
		q.OnError(err)
		return
	}
	if !ok {
		err = q.Remove(e.Notification.ID)
		if err != nil {
			q.OnError(err)
			return
		}
		q.OnRetryTooMany(e)
		return
	}
	e.ExecutionID, err = q.IDGenerator()
	err = q.Store.Replace(eid, e)
	if err != nil {
		q.OnError(err)
		return
	}
	q.pushExecution(e)
}
func (q *Queue) execute() {
	defer q.Recover()
	var iter = ""
	var err error
	var list []*notificationqueue.Execution
	for {
		list, iter, err = q.Store.List(iter, q.ExecuteCount)
		if err != nil {
			q.OnError(err)
			return
		}
		for _, e := range list {
			go q.retry(e)
		}
		if iter == "" {
			return
		}
	}
}
func (q *Queue) cron() {
	for {
		select {
		case <-time.After(q.Interval):
			go q.execute()
		case <-q.c:
			return
		}
	}
}
func (q *Queue) Start() error {
	q.c = make(chan int)
	go q.cron()
	return q.Store.Start()
}
func (q *Queue) Stop() error {
	close(q.c)
	return q.Store.Stop()
}

func New() *Queue {
	return &Queue{
		Interval:     DefaultInterval,
		ExecuteCount: DefaultExecuteCount,
		pipe:         make(chan *notificationqueue.Execution),
	}
}
