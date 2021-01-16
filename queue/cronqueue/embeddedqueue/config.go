package embeddedqueue

import (
	"github.com/herb-go/herbdata/kvdb"
	"github.com/herb-go/notification-drivers/queue/cronqueue"
	"github.com/herb-go/notification/notificationdelivery/notificationqueue"
)

type Config struct {
	Engine kvdb.Config
	Queue  cronqueue.Config
	Retry  *cronqueue.PlainRetry
}

var DirectiveName = "embeddedqueue"

func (c *Config) AppylToPublisher(p *notificationqueue.Publisher) error {
	q, err := c.Queue.CreateQueue()
	if err != nil {
		return err
	}
	db := kvdb.New()
	err = c.Engine.ApplyTo(db)
	if err != nil {
		return err
	}
	s := New()
	s.DB = db
	if c.Retry != nil && len(*c.Retry) > 0 {
		q.RetryHandler, err = c.Retry.CreateRetryHandler()
		if err != nil {
			return err
		}
	}
	err = q.AttachTo(p.Notifier)
	if err != nil {
		return err
	}
	q.Engine = s
	p.SetQueue(q)
	return nil
}

func Factory(loader func(v interface{}) error) (notificationqueue.Directive, error) {
	c := &Config{}
	err := loader(c)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func init() {
	notificationqueue.Register(DirectiveName, Factory)
}
