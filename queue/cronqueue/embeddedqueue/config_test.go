package embeddedqueue_test

import (
	"io/ioutil"
	"testing"

	"github.com/herb-go/herbdata-drivers/kvdb-drivers/leveldb"
	"github.com/herb-go/herbdata/kvdb"
	"github.com/herb-go/notification-drivers/queue/cronqueue"
	"github.com/herb-go/notification-drivers/queue/cronqueue/embeddedqueue"
	"github.com/herb-go/notification/notificationdelivery/notificationqueue"
)

func TestConfig(t *testing.T) {
	var err error
	tmpdir, err = ioutil.TempDir("", "")
	if err != nil {
		panic(err)
	}
	defer clean()
	c := &embeddedqueue.Config{
		Engine: kvdb.Config{
			Driver: "leveldb",
			Config: func(v interface{}) error {
				v.(*leveldb.Config).Database = tmpdir
				return nil
			},
		},
		Queue: cronqueue.Config{},
		Retry: &cronqueue.PlainRetry{"1s", "10s", "100s"},
	}
	d, err := notificationqueue.NewDirective(embeddedqueue.DirectiveName, func(v interface{}) error {
		v.(*embeddedqueue.Config).Engine = c.Engine
		v.(*embeddedqueue.Config).Queue = c.Queue
		v.(*embeddedqueue.Config).Retry = c.Retry
		return nil
	})
	if err != nil {
		panic(err)
	}
	p := notificationqueue.NewPublisher()
	err = d.AppylToPublisher(p)
	if err != nil {
		panic(err)
	}
	q, ok := p.Queue().(*cronqueue.Queue)
	if !ok {
		t.Fatal(ok)
	}
	if q.Engine == nil {
		t.Fatal(q)
	}
}
