package cronqueue

import (
	"io/ioutil"
	"os"
	"strconv"
	"sync/atomic"
	"testing"
	"time"

	"github.com/herb-go/herbdata-drivers/kvdb-drivers/leveldb"
	"github.com/herb-go/herbdata/kvdb"
	"github.com/herb-go/notification"
	"github.com/herb-go/notification-drivers/queue/cronqueue/embeddedstore"
	"github.com/herb-go/notification/notificationdelivery/notificationqueue"
)

var tmpdir string

func newTestStore() *embeddedstore.Store {
	s := embeddedstore.New()
	tmpdir, err := ioutil.TempDir("", "")
	if err != nil {
		panic(err)
	}
	s.DB = kvdb.New()
	c := leveldb.Config{
		Database: tmpdir,
	}
	d, err := c.CreateDriver()
	if err != nil {
		panic(err)
	}
	s.DB.Driver = d
	return s
}

func clean() {
	if tmpdir != "" {
		os.RemoveAll(tmpdir)
	}
}

var current int64

func idgen() (string, error) {
	c := atomic.AddInt64(&current, 1)
	return strconv.FormatInt(c, 10), nil
}

var errorlist = []error{}
var executionlist []*notificationqueue.Execution

func initTest() {
	errorlist = []error{}
	executionlist = []*notificationqueue.Execution{}
}
func testOnError(err error) {
	errorlist = append(errorlist, err)
}

func listen(c <-chan *notificationqueue.Execution) {
	go func() {
		for {
			select {
			case n, more := <-c:
				if !more {
					return
				}
				executionlist = append(executionlist, n)
			}
		}
	}()
}
func TestCronqueue(t *testing.T) {
	initTest()
	p := &PlainRetry{"15s", "12h"}
	r, err := p.CreateRetryHandler()
	if err != nil {
		t.Fatal(err)
	}
	q := New()
	q.OnError = testOnError
	q.Interval = time.Second
	q.IDGenerator = idgen
	q.RetryHandler = r
	q.Store = newTestStore()
	defer clean()
	c, err := q.PopChan()
	if err != nil {
		panic(err)
	}
	listen(c)
	err = q.Start()
	if err != nil {
		panic(err)
	}
	defer func() {
		err := q.Stop()
		if err != nil {
			panic(err)
		}
	}()
	n := notification.New()
	n.ID = "test"
	err = q.Push(n)
	if err != nil {
		panic(err)
	}
	time.Sleep(2 * time.Second)
}
