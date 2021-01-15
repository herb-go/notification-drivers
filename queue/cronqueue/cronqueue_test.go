package cronqueue_test

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
	"github.com/herb-go/notification-drivers/queue/cronqueue"
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
var timeoutlist []*notificationqueue.Execution
var retrytoomany []*notificationqueue.Execution

func initTest() {
	errorlist = []error{}
	executionlist = []*notificationqueue.Execution{}
	timeoutlist = []*notificationqueue.Execution{}
	retrytoomany = []*notificationqueue.Execution{}
}
func testOnError() {
	r := recover()
	if r != nil {
		err := r.(error)
		errorlist = append(errorlist, err)
	}
}

func testOnTimeout(e *notificationqueue.Execution) {
	timeoutlist = append(timeoutlist, e)
}
func testRetryTooMany(e *notificationqueue.Execution) {
	retrytoomany = append(retrytoomany, e)
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

func newTestQueue() *cronqueue.Queue {
	q := cronqueue.New()
	q.Interval = time.Second
	pr := &cronqueue.PlainRetry{"15s", "12h"}
	r, err := pr.CreateRetryHandler()
	if err != nil {
		panic(err)
	}
	n := notificationqueue.NewNotifier()
	n.Recover = testOnError
	n.OnDeliverTimeout = testOnTimeout
	n.OnRetryTooMany = testRetryTooMany
	n.IDGenerator = idgen
	err = q.AttachTo(n)
	if err != nil {
		panic(err)
	}
	q.RetryHandler = r
	q.Store = newTestStore()
	return q
}
func TestCronqueue(t *testing.T) {
	initTest()
	q := newTestQueue()
	q.Interval = time.Hour
	defer clean()
	p := notificationqueue.NewPublisher()
	p.IDGenerator = idgen
	err := q.AttachTo(p.Notifier)
	if err != nil {
		panic(err)
	}
	defer func() {
		err = q.Detach()
		if err != nil {
			panic(err)
		}
	}()
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
	if len(timeoutlist) != 0 {
		t.Fatal(len(timeoutlist))
	}
	if len(executionlist) != 1 {
		t.Fatal(len(executionlist))
	}
	if len(errorlist) != 0 {
		t.Fatal(len(errorlist))
	}
}
func TestTimeout(t *testing.T) {
	var err error
	initTest()
	q := newTestQueue()
	q.Timeout = time.Second
	q.Interval = time.Hour
	defer clean()
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
	if len(timeoutlist) != 0 {
		t.Fatal(len(timeoutlist))
	}
	n := notification.New()
	n.ID = "test"
	err = q.Push(n)
	if err != nil {
		panic(err)
	}
	time.Sleep(2 * time.Second)
	if len(timeoutlist) != 1 {
		t.Fatal(len(timeoutlist))
	}
	if len(errorlist) != 0 {
		t.Fatal(len(errorlist))
	}
}

func TestRetryTooMany(t *testing.T) {
	initTest()
	q := newTestQueue()
	p := &cronqueue.PlainRetry{"500ms", "500ms", "500ms"}
	r, err := p.CreateRetryHandler()
	if err != nil {
		panic(err)
	}
	q.Interval = 500 * time.Millisecond
	q.RetryHandler = r
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
	time.Sleep(5 * time.Second)
	if len(errorlist) != 0 {
		t.Fatal(len(errorlist))
	}
	if len(timeoutlist) != 0 {
		t.Fatal(len(timeoutlist))
	}
	if len(executionlist) != 3 {
		t.Fatal(len(executionlist))
	}
	if len(retrytoomany) != 1 {
		t.Fatal(len(retrytoomany))
	}

}
