package kvreceiptstore

import (
	"sync"
	"time"

	"github.com/herb-go/herbdata"
	"github.com/herb-go/herbdata/kvdb"
	"github.com/herb-go/notification"
	"github.com/herb-go/notification/notificationdelivery/notificationqueue"
	"github.com/vmihailenco/msgpack"
)

//New create new Store
func New() *Store {
	return &Store{}
}

type Store struct {
	locker               sync.Mutex
	DB                   *kvdb.Database
	DataRretentionPeriod time.Duration
	Limit                int
}

//Open open store and return any error if raised
func (s *Store) Open() error {
	return s.DB.Start()
}

//Close close store and return any error if raised
func (s *Store) Close() error {
	return s.DB.Stop()
}

//Save save given notificaiton to store.
//Receipt with same notification id will be overwritten.
func (s *Store) Save(receipt *notificationqueue.Receipt) error {
	s.locker.Lock()
	defer s.locker.Unlock()
	bs, err := msgpack.Marshal(receipt)
	if err != nil {
		return err
	}
	return s.DB.Set([]byte(receipt.Notification.ID), bs)

}

//List list no more than count notifactions in store with given search conditions form start position .
//Count should be greater than 0.
//Found receipts and next list position iter will be returned.
//Return largest id receipts if asc is false.
func (s *Store) List(condition []*notification.Condition, start string, asc bool, count int) (result []*notificationqueue.Receipt, iter string, err error) {
	var data []*herbdata.KeyValue
	var iterbs = []byte(start)
	var ok bool
	ctx := notification.NewConditionContext()
	limit := count
	if limit <= 0 {
		limit = notification.DefaultStoreListLimit
	}
	filter := NewFilter()
	err = ApplyToFilter(filter, condition)
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
			r := &notificationqueue.Receipt{}
			err = msgpack.Unmarshal(v.Value, r)
			if err != nil {
				return nil, "", err
			}
			ok, err = filter.FilterReceipt(r, ctx)
			if err != nil {
				return nil, "", err
			}
			if ok {
				result = append(result, r)
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

//Count count store with given search conditions
func (s *Store) Count(condition []*notification.Condition) (int, error) {
	var iter []byte
	var data []*herbdata.KeyValue
	var err error
	var result int
	var ok bool
	ctx := notification.NewConditionContext()
	limit := s.Limit
	if limit <= 0 {
		limit = notification.DefaultStoreListLimit
	}
	filter := NewFilter()
	err = ApplyToFilter(filter, condition)
	if err != nil {
		return 0, err
	}
	for {
		data, iter, err = s.DB.Next(iter, limit)
		if err != nil {
			return 0, err
		}
		for _, v := range data {
			r := &notificationqueue.Receipt{}
			err = msgpack.Unmarshal(v.Value, r)
			if err != nil {
				return 0, err
			}
			ok, err = filter.FilterReceipt(r, ctx)
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

//RetentionPeriod log retention period.
func (s *Store) RetentionPeriod() (time.Duration, error) {
	return s.DataRretentionPeriod, nil
}

//Clear clear outdate log
func (s *Store) Clear() error {
	return nil
}
