package embeddeddraftbox

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/herb-go/herbdata-drivers/kvdb-drivers/leveldb"
	"github.com/herb-go/herbdata/kvdb"
	"github.com/herb-go/notification-drivers/store/embeddedstore"
	"github.com/herb-go/notification/notificationdelivery/notificationqueue"
)

func TestDirective(t *testing.T) {
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
	dbconfig := &kvdb.Config{
		Driver: "leveldb",
		Config: func(v interface{}) error {
			lc := v.(*leveldb.Config)
			lc.Database = tmpdir
			return nil
		},
	}
	d, err := notificationqueue.NewDirective(DirectiveName, func(v interface{}) error {
		v.(*Config).Database = dbconfig
		return nil
	})
	if err != nil {
		panic(err)
	}
	p := notificationqueue.NewPublisher()
	d.AppylToPublisher(p)
	_, ok := p.Draftbox.(*embeddedstore.Store)
	if !ok {
		t.Fatal(d)
	}
}
