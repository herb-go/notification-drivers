package embeddedstore_test

import (
	"io/ioutil"
	"os"
	"testing"

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
func TestStore(t *testing.T) {
	s := newTestStore()
	defer clean()
	err := s.Start()
	if err != nil {
		panic(err)
	}
	defer func() {
		err := s.Stop()
		if err != nil {
			panic(err)
		}
	}()
	result, next, err := s.List("", 10)
	if len(result) != 0 || next != "" || err != nil {
		t.Fatal(result, next, err)
	}
	n := notification.New()
	n.ID = "test"
	e := notificationqueue.NewExecution()
	e.ExecutionID = "1"
	e.Notification = n
	err = s.Insert(e)
	if err != nil {
		t.Fatal(err)
	}
	result, next, err = s.List("", 10)
	if len(result) != 1 || next != "" || err != nil || result[0].ExecutionID != "1" {
		t.Fatal(result, next, err)
	}
	e.ExecutionID = "2"
	err = s.Insert(e)
	if err != nil {
		t.Fatal(err)
	}
	result, next, err = s.List("", 10)
	if len(result) != 1 || next != "" || err != nil || result[0].ExecutionID != "1" {
		t.Fatal(len(result), next, err)
	}
	err = s.Replace("notexist", e)
	if err != nil {
		t.Fatal(err)
	}
	result, next, err = s.List("", 10)
	if len(result) != 1 || next != "" || err != nil || result[0].ExecutionID != "1" {
		t.Fatal(len(result), next, err)
	}
	err = s.Replace("1", e)
	if err != nil {
		t.Fatal(err)
	}
	result, next, err = s.List("", 10)
	if len(result) != 1 || next != "" || err != nil || result[0].ExecutionID != "2" {
		t.Fatal(len(result), next, err)
	}
	err = s.Remove("notexist")
	if err != nil {
		t.Fatal(err)

	}
	result, next, err = s.List("", 10)
	if len(result) != 1 || next != "" || err != nil || result[0].ExecutionID != "2" {
		t.Fatal(len(result), next, err)
	}
	err = s.Remove(n.ID)
	if err != nil {
		t.Fatal(err)

	}
	result, next, err = s.List("", 10)
	if len(result) != 0 || next != "" || err != nil {
		t.Fatal(len(result), next, err)
	}
	n.ID = "1"
	err = s.Insert(e)
	if err != nil {
		t.Fatal(err)
	}
	n.ID = "2"
	err = s.Insert(e)
	if err != nil {
		t.Fatal(err)
	}
	n.ID = "3"
	err = s.Insert(e)
	if err != nil {
		t.Fatal(err)
	}
	result, next, err = s.List("", 2)
	if len(result) != 2 || next == "" || err != nil {
		t.Fatal(len(result), next, err)
	}
	result, next, err = s.List(next, 2)
	if len(result) != 1 || next != "" || err != nil {
		t.Fatal(len(result), next, err)
	}
}
