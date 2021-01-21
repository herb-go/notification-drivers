package passthroughqueue

import (
	"time"

	"github.com/herb-go/notification"
	"github.com/herb-go/notification/notificationdelivery/notificationqueue"
)

var DirectiveName = "passthroughqueue"

type Passthroughqueue struct {
	c           chan *notificationqueue.Execution
	IDGenerator func() (string, error)
}

//PopChan return execution chan
func (q *Passthroughqueue) PopChan() (<-chan *notificationqueue.Execution, error) {
	return q.c, nil
}

//Push push notification to queue
func (q *Passthroughqueue) Push(n *notification.Notification) error {
	var err error
	e := notificationqueue.NewExecution()
	e.ExecutionID, err = q.IDGenerator()
	if err != nil {
		return err
	}
	e.StartTime = time.Now().Unix()
	e.Notification = n
	go func() {
		q.c <- e
	}()
	return nil
}

//Remove remove notification by given id
func (q *Passthroughqueue) Remove(nid string) error {
	return nil
}

//Start start queue
func (q *Passthroughqueue) Start() error {
	return nil
}

//Stop stop queue
func (q *Passthroughqueue) Stop() error {
	close(q.c)
	return nil
}

//AttachTo attach queue to notifier
func (q *Passthroughqueue) AttachTo(n *notificationqueue.Notifier) error {
	q.IDGenerator = func() (string, error) {
		return n.IDGenerator()
	}
	return nil
}

//Detach detach queue.
func (pq *Passthroughqueue) Detach() error {
	return nil
}

func (q *Passthroughqueue) AppylToPublisher(p *notificationqueue.Publisher) error {

	p.SetQueue(q)
	return nil
}
func New() *Passthroughqueue {
	return &Passthroughqueue{
		c: make(chan *notificationqueue.Execution),
	}
}
func Factory(loader func(v interface{}) error) (notificationqueue.Directive, error) {
	return New(), nil
}

func init() {
	notificationqueue.Register(DirectiveName, Factory)
}
