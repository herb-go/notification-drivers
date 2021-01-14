package embeddedstore

import (
	"sync"

	"github.com/vmihailenco/msgpack"

	"github.com/herb-go/herbdata"

	"github.com/herb-go/herbdata/kvdb"
	"github.com/herb-go/notification/notificationdelivery/notificationqueue"
)

type Store struct {
	locker sync.Mutex
	DB     *kvdb.Database
}

func (s *Store) List(start string, count int) ([]*notificationqueue.Execution, string, error) {
	s.locker.Lock()
	defer s.locker.Unlock()
	iter := []byte(start)
	data, iter, err := s.DB.Next(iter, count)
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
func (s *Store) Insert(execution *notificationqueue.Execution) error {
	s.locker.Lock()
	defer s.locker.Unlock()
	bs, err := msgpack.Marshal(execution)
	if err != nil {
		return err
	}
	_, err = s.DB.Get([]byte(execution.Notification.ID))
	if err == nil {
		return nil
	}
	if err != herbdata.ErrNotFound {
		return err
	}
	return s.DB.Set([]byte(execution.Notification.ID), bs)
}
func (s *Store) Replace(eid string, new *notificationqueue.Execution) error {
	s.locker.Lock()
	defer s.locker.Unlock()
	bs, err := s.DB.Get([]byte(new.Notification.ID))
	if err != nil {
		return err
	}
	e := notificationqueue.NewExecution()
	err = msgpack.Unmarshal(bs, e)
	if err != nil {
		return err
	}
	if e.ExecutionID != eid {
		return nil
	}
	bs, err = msgpack.Marshal(new)
	if err != nil {
		return err
	}
	return s.DB.Set([]byte(new.Notification.ID), bs)
}
func (s *Store) Remove(nid string) error {
	s.locker.Lock()
	defer s.locker.Unlock()
	return s.DB.Delete([]byte(nid))
}
func (s *Store) Start() error {
	return s.DB.Start()
}
func (s *Store) Stop() error {
	return s.DB.Stop()
}

func New() *Store {
	return &Store{}
}
