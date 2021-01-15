package cronqueue_test

import (
	"testing"
	"time"

	"github.com/herb-go/notification-drivers/queue/cronqueue"
)

func TestConfig(t *testing.T) {
	c := &cronqueue.Config{}
	c.TimeoutDuration = "10s"
	c.IntervalDuration = "15s"
	c.ExecuteCount = 999
	q, err := c.CreateQueue()
	if err != nil {
		panic(err)
	}
	if q.Timeout != 10*time.Second || q.Interval != 15*time.Second || q.ExecuteCount != 999 {
		t.Fatal(q)
	}
}
