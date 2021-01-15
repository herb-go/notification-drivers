package embeddeddraftbox

import (
	"sync"
	"time"

	"github.com/herb-go/herbdata"

	"github.com/herb-go/herbdata/kvdb"

	"github.com/herb-go/notification"
	"github.com/herb-go/notification/notificationdelivery/notificationqueue"
	"github.com/vmihailenco/msgpack"
)

//DirectiveName registered direcvive name
const DirectiveName = "embeddeddraftbox"

//RequiredKvdbFeatures required kvdb featuers
var RequiredKvdbFeatures = kvdb.FeatureEmbedded | kvdb.FeatureStore | kvdb.FeatureNext | kvdb.FeaturePrev

//SupportedConditions supported filter conditions
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

//Draftbox draftbox struct
type Draftbox struct {
	locker sync.Mutex
	DB     *kvdb.Database
	Limit  int
}

//New create new draft box
func New() *Draftbox {
	return &Draftbox{}
}

//Open open draftbox and return any error if raised
func (d *Draftbox) Open() error {
	return d.DB.Start()
}

//Close close draftbox and return any error if raised
func (d *Draftbox) Close() error {
	return d.DB.Stop()
}

//Draft save given notificaiton to draft box.
//Notification with same id will be overwritten.
func (d *Draftbox) Draft(n *notification.Notification) error {
	d.locker.Lock()
	defer d.locker.Unlock()
	bs, err := msgpack.Marshal(n)
	if err != nil {
		return err
	}
	return d.DB.Set([]byte(n.ID), bs)
}

//List list no more than count notifactions in draftbox with given search conditions form start position .
//Count should be greater than 0.
//Found notifications and next list position iter will be returned.
//Return largest id notification if asc is false.
func (d *Draftbox) List(condition []*notificationqueue.Condition, start string, asc bool, count int) (result []*notification.Notification, iter string, err error) {
	var data []*herbdata.KeyValue
	var iterbs = []byte(start)
	var ok bool
	ts := time.Now().Unix()
	limit := count
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
			data, iterbs, err = d.DB.Next(iterbs, limit)
		} else {
			data, iterbs, err = d.DB.Prev(iterbs, limit)
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

//Count draft box with given search conditions
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

//SupportedConditions return supported condition keyword list
func (d *Draftbox) SupportedConditions() ([]string, error) {
	return SupportedConditions, nil
}

//Eject remove notification by given id and return removed notification.
func (d *Draftbox) Eject(id string) (*notification.Notification, error) {
	d.locker.Lock()
	defer d.locker.Unlock()
	bs, err := d.DB.Get([]byte(id))
	if err != nil {
		if err == herbdata.ErrNotFound {
			return nil, notification.NewErrNotificationIDNotFound(id)
		}
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

//Config draft box config
type Config struct {
	//Database kvdb config
	Database *kvdb.Config
	//Limit count limit,defalut value is notificationquque.DefaultDraftboxListLimit
	Limit int
}

//CreateDraftbox create draftbox with config
func (c *Config) CreateDraftbox() (notificationqueue.Draftbox, error) {
	var err error
	d := New()
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

//AppylToPublisher create draft box and add to publisher
func (c *Config) AppylToPublisher(p *notificationqueue.Publisher) error {
	draftbox, err := c.CreateDraftbox()
	if err != nil {
		return err
	}
	p.Draftbox = draftbox
	return nil
}

//Factory draftbox directive factory
var Factory = func(loader func(v interface{}) error) (notificationqueue.Directive, error) {
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
