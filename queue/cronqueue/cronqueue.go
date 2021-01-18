package cronqueue

import (
	"time"

	"github.com/herb-go/notification"
	"github.com/herb-go/notification/notificationdelivery/notificationqueue"
)

//DefaultInterval default queue interval
var DefaultInterval = time.Minute

//DefaultExecuteCount default execute count
var DefaultExecuteCount = 10

//DefaultTimeout default push timeout
var DefaultTimeout = time.Minute

//Queue queue struct
type Queue struct {
	Engine           Engine
	Timeout          time.Duration
	Interval         time.Duration
	ExecuteCount     int
	IDGenerator      func() (string, error)
	c                chan int
	pipe             chan *notificationqueue.Execution
	Recover          func()
	OnExecution      func(*notificationqueue.Execution)
	OnDeliverTimeout func(*notificationqueue.Execution)
	OnRetryTooMany   func(*notificationqueue.Execution)
	RetryHandler     RetryHandler
}

//PopChan return execution chan
func (q *Queue) PopChan() (<-chan *notificationqueue.Execution, error) {
	return q.pipe, nil
}

func (q *Queue) createExecution(n *notification.Notification) (*notificationqueue.Execution, error) {
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
	go q.OnExecution(e)
	select {
	case q.pipe <- e:
	case <-time.After(q.Timeout):
		go q.OnDeliverTimeout(e)
	}
}

//Push push notification to queue
func (q *Queue) Push(n *notification.Notification) error {
	e, err := q.createExecution(n)
	if err != nil {
		return err
	}
	ok, err := q.RetryHandler.HandleRetry(e)
	if err != nil {
		return err
	}
	if ok {
		err = q.Engine.Insert(e)
		if err != nil {
			return err
		}
	}
	go q.pushExecution(e)
	return nil
}

//Remove remove notification by given id
func (q *Queue) Remove(nid string) error {
	return q.Engine.Remove(nid)
}

func (q *Queue) retry(e *notificationqueue.Execution) {
	defer q.Recover()
	if time.Now().Unix() <= e.RetryAfterTime {
		return
	}
	eid := e.ExecutionID
	ok, err := q.RetryHandler.HandleRetry(e)
	if err != nil {
		panic(err)
	}
	if !ok {
		err = q.Remove(e.Notification.ID)
		if err != nil {
			panic(err)
		}
		q.OnRetryTooMany(e)
		return
	}
	e.ExecutionID, err = q.IDGenerator()
	err = q.Engine.Replace(eid, e)
	if err != nil {
		panic(err)
	}
	q.pushExecution(e)
}
func (q *Queue) execute() {
	defer q.Recover()
	var iter = ""
	var err error
	var list []*notificationqueue.Execution
	for {
		list, iter, err = q.Engine.List(iter, q.ExecuteCount)
		if err != nil {
			panic(err)
		}
		for _, e := range list {
			q.retry(e)
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

//Start start queue
func (q *Queue) Start() error {
	q.c = make(chan int)
	go q.cron()
	return q.Engine.Start()
}

//Stop stop queue
func (q *Queue) Stop() error {
	close(q.c)
	return q.Engine.Stop()
}

//AttachTo attach queue to notifier
func (q *Queue) AttachTo(n *notificationqueue.Notifier) error {
	q.OnDeliverTimeout = func(e *notificationqueue.Execution) {
		n.OnDeliverTimeout(e)
	}
	q.OnRetryTooMany = func(e *notificationqueue.Execution) {
		n.OnRetryTooMany(e)
	}
	q.Recover = func() {
		defer n.Recover()
		if r := recover(); r != nil {
			panic(r)
		}
	}
	q.IDGenerator = func() (string, error) {
		return n.IDGenerator()
	}
	q.OnExecution = func(e *notificationqueue.Execution) {
		n.OnExecution(e)
	}
	return nil
}

//Detach detach queue.
func (q *Queue) Detach() error {
	q.OnDeliverTimeout = notificationqueue.NopExecutionHandler
	q.OnRetryTooMany = notificationqueue.NopExecutionHandler
	q.Recover = notificationqueue.NopRecover
	q.IDGenerator = notificationqueue.NopIDGenerator
	return nil
}

//New create new queue
func New() *Queue {
	return &Queue{
		Interval:     DefaultInterval,
		Timeout:      DefaultTimeout,
		ExecuteCount: DefaultExecuteCount,
		pipe:         make(chan *notificationqueue.Execution),
	}
}
