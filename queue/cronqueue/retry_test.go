package cronqueue

import (
	"testing"

	"github.com/herb-go/notification"
	"github.com/herb-go/notification/notificationdelivery/notificationqueue"
)

func TestRetry(t *testing.T) {
	var ok bool
	p := PlainRetry{"15s", "12h"}
	r, err := p.CreateRetryHandler()
	if len(*r) != 2 || err != nil {
		t.Fatal(r, err)
	}
	e := notificationqueue.NewExecution()
	e.Notification = notification.New()
	ok, err = r.HandleRetry(e)
	if !ok || err != nil {
		t.Fatal(ok, err)
	}
	if e.RetryCount != 1 || e.RetryAfterTime-e.StartTime != 15 {
		t.Fatal(e)
	}
	ok, err = r.HandleRetry(e)
	if !ok || err != nil {
		t.Fatal(ok, err)
	}
	if e.RetryCount != 2 || e.RetryAfterTime-e.StartTime != 12*3600 {
		t.Fatal(e)
	}
	ok, err = r.HandleRetry(e)
	if ok || err != nil {
		t.Fatal(ok, err)
	}
}
