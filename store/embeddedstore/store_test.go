package embeddedstore_test

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/herb-go/herbdata-drivers/kvdb-drivers/leveldb"
	_ "github.com/herb-go/herbdata-drivers/kvdb-drivers/leveldb"
	"github.com/herb-go/herbdata/kvdb"
	"github.com/herb-go/notification"
	"github.com/herb-go/notification-drivers/store/embeddedstore"
)

func TestStore(t *testing.T) {
	tmpdir, err := ioutil.TempDir("", "")
	if err != nil {
		panic(err)
	}
	defer func() {
		if tmpdir == "" {
			return
		}
		os.RemoveAll(tmpdir)
	}()
	c := &embeddedstore.Config{
		Database: &kvdb.Config{
			Driver: "leveldb",
			Config: func(v interface{}) error {
				lc := v.(*leveldb.Config)
				lc.Database = tmpdir
				return nil
			},
		},
	}
	d, err := c.CreateStore()
	if err != nil {
		panic(err)
	}
	err = d.Open()
	if err != nil {
		panic(err)
	}
	defer func() {
		err = d.Close()
		if err != nil {
			panic(err)
		}
	}()
	dbox := d.(*embeddedstore.Store)
	if dbox.Limit != notification.DefaultStoreListLimit {
		t.Fatal(dbox)
	}
	dbox.Limit = 0
	supported, err := d.SupportedConditions()
	if err != nil || strings.Join(supported, ",") != strings.Join(notification.PlainFilterSupportedConditions, ",") {
		t.Fatal(supported, err)
	}
	n := notification.New()
	n.ID = "1"
	err = d.Save(n)
	if err != nil {
		panic(err)
	}
	cont, err := d.Count(nil)
	if err != nil || cont != 1 {
		t.Fatal(cont, err)
	}
	n = notification.New()
	n.ID = "2"
	err = d.Save(n)
	if err != nil {
		panic(err)
	}
	cont, err = d.Count(nil)
	if err != nil || cont != 2 {
		t.Fatal(cont, err)
	}
	n = notification.New()
	n.ID = "2"
	err = d.Save(n)
	if err != nil {
		panic(err)
	}
	cont, err = d.Count(nil)
	if err != nil || cont != 2 {
		t.Fatal(cont, err)
	}
	conds := []*notification.Condition{&notification.Condition{Keyword: notification.ConditionNotificationID, Value: "1"}}
	cont, err = d.Count(conds)
	if err != nil || cont != 1 {
		t.Fatal(cont, err)
	}
	n, err = d.Remove("1")
	if err != nil || n == nil || n.ID != "1" {
		t.Fatal(n, err)
	}
	n, err = d.Remove("2")
	if err != nil || n == nil || n.ID != "2" {
		t.Fatal(n, err)
	}
	cont, err = d.Count(nil)
	if err != nil || cont != 0 {
		t.Fatal(cont, err)
	}
	n, err = d.Remove("2")
	if !notification.IsErrNotificationIDNotFound(err) {
		t.Fatal(err)
	}
	n = notification.New()
	n.ID = "1"
	err = d.Save(n)
	if err != nil {
		t.Fatal(err)
	}
	n = notification.New()
	n.ID = "2"
	n.Header.Set(notification.HeaderNameBatch, "batch")
	err = d.Save(n)
	if err != nil {
		t.Fatal(err)
	}
	n = notification.New()
	n.ID = "3"
	err = d.Save(n)
	if err != nil {
		t.Fatal(err)
	}
	n = notification.New()
	n.ID = "4"
	n.Header.Set(notification.HeaderNameBatch, "batch")
	err = d.Save(n)
	if err != nil {
		t.Fatal(err)
	}
	n = notification.New()
	n.ID = "5"
	err = d.Save(n)
	if err != nil {
		t.Fatal(err)
	}
	n = notification.New()
	n.ID = "6"
	n.Header.Set(notification.HeaderNameBatch, "batch")
	err = d.Save(n)
	if err != nil {
		t.Fatal(err)
	}
	cont, err = d.Count(nil)
	if err != nil || cont != 6 {
		t.Fatal(cont, err)
	}
	result, iter, err := d.List(nil, "", false, 0)
	if len(result) != 6 || iter != "" || err != nil {
		t.Fatal(len(result), iter, err)
	}
	if result[0].ID != "6" || result[1].ID != "5" || result[2].ID != "4" || result[3].ID != "3" || result[4].ID != "2" || result[5].ID != "1" {
		t.Fatal(result)
	}

	result, iter, err = d.List(nil, "", false, 4)
	if len(result) != 4 || iter == "" || err != nil {
		t.Fatal(len(result), iter, err)
	}
	if result[0].ID != "6" || result[1].ID != "5" || result[2].ID != "4" || result[3].ID != "3" {
		t.Fatal(result)
	}
	result, iter, err = d.List(nil, iter, false, 4)
	if len(result) != 2 || iter != "" || err != nil {
		t.Fatal(len(result), iter, err)
	}
	if result[0].ID != "2" || result[1].ID != "1" {
		t.Fatal(result)
	}

	result, iter, err = d.List(nil, "", true, 0)
	if len(result) != 6 || iter != "" || err != nil {
		t.Fatal(len(result), iter, err)
	}
	if result[0].ID != "1" || result[1].ID != "2" || result[2].ID != "3" || result[3].ID != "4" || result[4].ID != "5" || result[5].ID != "6" {
		t.Fatal(result)
	}

	result, iter, err = d.List(nil, "", true, 4)
	if len(result) != 4 || iter == "" || err != nil {
		t.Fatal(len(result), iter, err)
	}
	if result[0].ID != "1" || result[1].ID != "2" || result[2].ID != "3" || result[3].ID != "4" {
		t.Fatal(result)
	}
	result, iter, err = d.List(nil, iter, true, 4)
	if len(result) != 2 || iter != "" || err != nil {
		t.Fatal(len(result), iter, err)
	}
	if result[0].ID != "5" || result[1].ID != "6" {
		t.Fatal(result)
	}

	conds = []*notification.Condition{&notification.Condition{Keyword: notification.ConditionBatch, Value: "batch"}}
	cont, err = d.Count(conds)
	if err != nil || cont != 3 {
		t.Fatal(cont, err)
	}

	result, iter, err = d.List(conds, "", false, 0)
	if len(result) != 3 || iter != "" || err != nil {
		t.Fatal(len(result), iter, err)
	}
	if result[0].ID != "6" || result[1].ID != "4" || result[2].ID != "2" {
		t.Fatal(result)
	}

	result, iter, err = d.List(conds, "", false, 2)
	if len(result) != 2 || iter == "" || err != nil {
		t.Fatal(len(result), iter, err)
	}
	if result[0].ID != "6" || result[1].ID != "4" {
		t.Fatal(result)
	}
	result, iter, err = d.List(conds, iter, false, 2)
	if len(result) != 1 || iter != "" || err != nil {
		t.Fatal(len(result), iter, err)
	}
	if result[0].ID != "2" {
		t.Fatal(result)
	}

	result, iter, err = d.List(conds, "", true, 0)
	if len(result) != 3 || iter != "" || err != nil {
		t.Fatal(len(result), iter, err)
	}
	if result[0].ID != "2" || result[1].ID != "4" || result[2].ID != "6" {
		t.Fatal(result)
	}

	result, iter, err = d.List(conds, "", true, 2)
	if len(result) != 2 || iter == "" || err != nil {
		t.Fatal(len(result), iter, err)
	}
	if result[0].ID != "2" || result[1].ID != "4" {
		t.Fatal(result)
	}
	result, iter, err = d.List(conds, iter, true, 2)
	if len(result) != 1 || iter != "" || err != nil {
		t.Fatal(len(result), iter, err)
	}
	if result[0].ID != "6" {
		t.Fatal(result)
	}
}
