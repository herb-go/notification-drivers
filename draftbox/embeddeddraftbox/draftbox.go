package embeddeddraftbox

import (
	"sync"
	"time"

	"github.com/herb-go/herbdata"

	"github.com/herb-go/herbdata/kvdb"

	"github.com/herb-go/notification"
	"github.com/herb-go/notification/notificationqueue"
	"github.com/vmihailenco/msgpack"
)

//RequiredKvdbFeatures required kvdb featuers
var RequiredKvdbFeatures = kvdb.FeatureEmbedded | kvdb.FeatureStore | kvdb.FeatureNext | kvdb.FeaturePrev

var SupportedConditions = []string{
	notificationqueue.ConditionBatch,
	notificationqueue.ConditionNotificationID,
	notificationqueue.ConditionDelivery,
	notificationqueue.ConditionTarget,
	notificationqueue.ConditionTopic,
	notificationqueue.ConditionInContent,
	notificationqueue.ConditionBeforeTimestamp,
	notificationqueue.ConditionAfterTimestamp,
}

type Draftbox struct {
	locker sync.Mutex
	DB     *kvdb.Database
	Limit  int
}

func New() *Draftbox {
	return &Draftbox{}
}
func (d *Draftbox) Open() error {
	return d.DB.Start()
}
func (d *Draftbox) Close() error {
	return d.DB.Stop()
}
func (d *Draftbox) Draft(n *notification.Notification) error {
	d.locker.Lock()
	defer d.locker.Unlock()
	bs, err := msgpack.Marshal(n)
	if err != nil {
		return err
	}
	return d.DB.Set([]byte(n.ID), bs)
}
func (d *Draftbox) List(condition []*notificationqueue.Condition, start string, asc bool, count int) (result []*notification.Notification, iter string, err error) {
	var data []*herbdata.KeyValue
	var iterbs []byte
	var ok bool
	ts := time.Now().Unix()
	limit := d.Limit
	if limit <= 0 {
		limit = notificationqueue.DefaultDraftboxListLimit
	}
	filter := notificationqueue.NewFilter()
	err = notificationqueue.ApplyToFilter(filter, condition)
	if err != nil {
		return nil, "", err
	}
	for {
		if asc {
			data, iterbs, err = d.DB.Prev(iterbs, limit)
		} else {
			data, iterbs, err = d.DB.Next(iterbs, limit)
		}
		if err != nil {
			return nil, "", err
		}
		for _, v := range data {
			n := &notification.Notification{}
			err = msgpack.Unmarshal(v.Value, n)
			if err != nil {
				return nil, "", err
			}
			ok, err = filter.FilterNotification(n, ts)
			if err != nil {
				return nil, "", err
			}
			if ok {
				result = append(result, n)
			}

		}
		if len(iterbs) == 0 {
			break
		}
		if len(result) >= limit {
			return result, string(iterbs), nil
		}
	}
	return result, "", nil
}
func (d *Draftbox) Count(condition []*notificationqueue.Condition) (int, error) {
	var iter []byte
	var data []*herbdata.KeyValue
	var err error
	var result int
	var ok bool
	ts := time.Now().Unix()
	limit := d.Limit
	if limit <= 0 {
		limit = notificationqueue.DefaultDraftboxListLimit
	}
	filter := notificationqueue.NewFilter()
	err = notificationqueue.ApplyToFilter(filter, condition)
	if err != nil {
		return 0, err
	}
	for {
		data, iter, err = d.DB.Next(iter, limit)
		if err != nil {
			return 0, err
		}
		for _, v := range data {
			n := &notification.Notification{}
			err = msgpack.Unmarshal(v.Value, n)
			if err != nil {
				return 0, err
			}
			ok, err = filter.FilterNotification(n, ts)
			if err != nil {
				return 0, err
			}
			if ok {
				result = result + 1
			}
		}
		if len(iter) == 0 {
			break
		}
	}
	return result, nil
}
func (d *Draftbox) SupportedConditions() ([]string, error) {
	return SupportedConditions, nil
}
func (d *Draftbox) Eject(id string) (*notification.Notification, error) {
	d.locker.Lock()
	defer d.locker.Unlock()
	bs, err := d.DB.Get([]byte(id))
	if err != nil {
		return nil, err
	}
	err = d.DB.Delete([]byte(id))
	if err != nil {
		return nil, err
	}
	n := notification.New()
	err = msgpack.Unmarshal(bs, n)
	if err != nil {
		return nil, err
	}
	return n, nil
}

type Config struct {
	Database *kvdb.Config
	Limit    int
}

func (c *Config) CreateDraftbox() (notificationqueue.Draftbox, error) {
	var err error
	d := &Draftbox{}
	d.DB = kvdb.New()
	err = c.Database.ApplyTo(d.DB)
	if err != nil {
		return nil, err
	}
	d.Limit = c.Limit
	if d.Limit <= 0 {
		d.Limit = notificationqueue.DefaultDraftboxListLimit
	}
	err = d.DB.ShouldSupport(RequiredKvdbFeatures)
	if err != nil {
		return nil, err
	}
	return d, nil
}
