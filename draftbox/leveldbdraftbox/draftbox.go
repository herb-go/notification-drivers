package leveldbdraftbox

import (
	"strconv"
	"strings"
	"sync"

	"github.com/herb-go/notification"
	"github.com/herb-go/notification/notificationqueue"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/vmihailenco/msgpack"
)

var Supported = []string{
	notificationqueue.ConditionBatch,
	notificationqueue.ConditionNotificationID,
	notificationqueue.ConditionDelivery,
	notificationqueue.ConditionTarget,
	notificationqueue.ConditionTopic,
	notificationqueue.ConditionInContent,
	notificationqueue.ConditionBeforeTimestamp,
	notificationqueue.ConditionAfterTimestamp,
}

type Condition struct {
	BatchID        string
	NotificationID string
	Delivery       string
	Target         string
	Topic          string
	InContent      string
	After          int64
	Before         int64
}

func (c *Condition) applyCondition(cond *notificationqueue.Condition) error {
	switch cond.Keyword {
	case notificationqueue.ConditionBatch:
		c.BatchID = cond.Value
	case notificationqueue.ConditionNotificationID:
		c.NotificationID = cond.Value
	case notificationqueue.ConditionDelivery:
		c.Delivery = cond.Value
	case notificationqueue.ConditionTarget:
		c.Target = cond.Value
	case notificationqueue.ConditionTopic:
		c.Topic = cond.Value
	case notificationqueue.ConditionAfterTimestamp:
		ts, err := strconv.ParseInt(cond.Value, 10, 64)
		if err != nil {
			return nil
		}
		c.After = ts
	case notificationqueue.ConditionBeforeTimestamp:
		ts, err := strconv.ParseInt(cond.Value, 10, 64)
		if err != nil {
			return nil
		}
		c.Before = ts
	default:
		return notificationqueue.ErrConditionNotSupported(cond.Value)
	}
	return nil
}
func (c *Condition) Apply(conds []*notificationqueue.Condition) error {
	for k := range conds {
		err := c.applyCondition(conds[k])
		if err != nil {
			return err
		}
	}
	return nil
}
func (c *Condition) Filter(n *notification.Notification, ts int64) bool {
	if c.BatchID != "" && n.Header.Get(notification.HeaderNameBatch) != c.BatchID {
		return false
	}
	if c.NotificationID != "" && n.ID != c.NotificationID {
		return false
	}
	if c.Delivery != "" && n.Delivery != c.Delivery {
		return false
	}
	if c.Target != "" && n.Header.Get(notification.HeaderNameTarget) != c.Target {
		return false
	}
	if c.Topic != "" && n.Header.Get(notification.HeaderNameTopic) != c.Topic {
		return false
	}
	if c.InContent != "" {
		var found bool
		for k := range n.Content {
			if strings.Contains(n.Content[k], c.InContent) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	if c.After > 0 && c.After <= ts {
		return false
	}
	if c.Before > 0 && c.After >= ts {
		return false
	}
	return true
}

type Draftbox struct {
	locker sync.Mutex
	leveldb.DB
}

func (d *Draftbox) Draft(n *notification.Notification) error {
	d.locker.Lock()
	defer d.locker.Unlock()
	bs, err := msgpack.Marshal(n)
	if err != nil {
		return err
	}
	return d.DB.Put([]byte(n.ID), bs, nil)
}
func (d *Draftbox) List(condition []*notificationqueue.Condition, start string, asc bool, count int) (result []*notification.Notification, iter string, err error) {

}
func (d *Draftbox) Count(condition []*notificationqueue.Condition) (int, error) {

}
func (d *Draftbox) SupportedConditions() ([]string, error) {
	return Supported, nil
}
func (d *Draftbox) Eject(id string) (*notification.Notification, error) {
	d.locker.Lock()
	defer d.locker.Unlock()
	bs, err := d.DB.Get([]byte(id), nil)
	if err != nil {
		return nil, err
	}
	n := notification.New()
	err = msgpack.Unmarshal(bs, n)
	if err != nil {
		return nil, err
	}
	err = d.DB.Delete([]byte(id), nil)
	if err != nil {
		return nil, err
	}

	return n, nil
}
