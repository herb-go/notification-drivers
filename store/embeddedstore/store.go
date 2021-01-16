package embeddedstore

import (
	"sync"
	"time"

	"github.com/herb-go/herbdata"

	"github.com/herb-go/herbdata/kvdb"

	"github.com/herb-go/notification"
	"github.com/vmihailenco/msgpack"
)

//RequiredKvdbFeatures required kvdb featuers
var RequiredKvdbFeatures = kvdb.FeatureEmbedded | kvdb.FeatureStore | kvdb.FeatureNext | kvdb.FeaturePrev

//SupportedConditions supported filter conditions
var SupportedConditions = []string{
	notification.ConditionBatch,
	notification.ConditionNotificationID,
	notification.ConditionDelivery,
	notification.ConditionTarget,
	notification.ConditionTopic,
	notification.ConditionInContent,
	notification.ConditionBeforeTimestamp,
	notification.ConditionAfterTimestamp,
}

//Store draftbox struct
type Store struct {
	locker sync.Mutex
	DB     *kvdb.Database
	Limit  int
}

//New create new draft box
func New() *Store {
	return &Store{}
}

//Open open draftbox and return any error if raised
func (s *Store) Open() error {
	return s.DB.Start()
}

//Close close draftbox and return any error if raised
func (s *Store) Close() error {
	return s.DB.Stop()
}

//Save save given notificaiton to draft box.
//Notification with same id will be overwritten.
func (s *Store) Save(n *notification.Notification) error {
	s.locker.Lock()
	defer s.locker.Unlock()
	bs, err := msgpack.Marshal(n)
	if err != nil {
		return err
	}
	return s.DB.Set([]byte(n.ID), bs)
}

//List list no more than count notifactions in draftbox with given search conditions form start position .
//Count should be greater than 0.
//Found notifications and next list position iter will be returnes.
//Return largest id notification if asc is false.
func (s *Store) List(condition []*notification.Condition, start string, asc bool, count int) (result []*notification.Notification, iter string, err error) {
	var data []*herbdata.KeyValue
	var iterbs = []byte(start)
	var ok bool
	ts := time.Now().Unix()
	limit := count
	if limit <= 0 {
		limit = notification.DefaultStoreListLimit
	}
	filter := notification.NewFilter()
	err = notification.ApplyToFilter(filter, condition)
	if err != nil {
		return nil, "", err
	}
	for {
		if asc {
			data, iterbs, err = s.DB.Next(iterbs, limit)
		} else {
			data, iterbs, err = s.DB.Prev(iterbs, limit)
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
func (s *Store) Count(condition []*notification.Condition) (int, error) {
	var iter []byte
	var data []*herbdata.KeyValue
	var err error
	var result int
	var ok bool
	ts := time.Now().Unix()
	limit := s.Limit
	if limit <= 0 {
		limit = notification.DefaultStoreListLimit
	}
	filter := notification.NewFilter()
	err = notification.ApplyToFilter(filter, condition)
	if err != nil {
		return 0, err
	}
	for {
		data, iter, err = s.DB.Next(iter, limit)
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
func (s *Store) SupportedConditions() ([]string, error) {
	return SupportedConditions, nil
}

//Remove remove notification by given id and return removed notification.
func (s *Store) Remove(id string) (*notification.Notification, error) {
	s.locker.Lock()
	defer s.locker.Unlock()
	bs, err := s.DB.Get([]byte(id))
	if err != nil {
		if err == herbdata.ErrNotFound {
			return nil, notification.NewErrNotificationIDNotFound(id)
		}
		return nil, err
	}
	err = s.DB.Delete([]byte(id))
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
	//Limit count limit,defalut value is notificationquque.DefaultStoreListLimit
	Limit int
}

//CreateStore create draftbox with config
func (c *Config) CreateStore() (notification.Store, error) {
	var err error
	s := New()
	s.DB = kvdb.New()
	err = c.Database.ApplyTo(s.DB)
	if err != nil {
		return nil, err
	}
	s.Limit = c.Limit
	if s.Limit <= 0 {
		s.Limit = notification.DefaultStoreListLimit
	}
	err = s.DB.ShouldSupport(RequiredKvdbFeatures)
	if err != nil {
		return nil, err
	}
	return s, nil
}
