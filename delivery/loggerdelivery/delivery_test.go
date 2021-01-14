package loggerdelivery

import (
	"bytes"
	"strings"
	"testing"

	"github.com/herb-go/logger"

	"github.com/herb-go/notification"
	"github.com/herb-go/notification/notificationdelivery"
)

func newAddr(name string, mail string) string {
	d := &Delivery{}
	return d.MustEscape(name) + " <" + d.MustEscape(mail) + ">"
}
func NewTestDelivery() *Delivery {
	dc := &notificationdelivery.Config{
		DeliveryType: DeliveryType,
		DeliveryConfig: func(v interface{}) error {
			v.(*Config).Logger = "debug"
			return nil
		},
	}
	d, err := dc.CreateDriver()
	if err != nil {
		panic(err)
	}
	return d.(*Delivery)
}

var _ notificationdelivery.DeliveryDriver = &Delivery{}

func TestDelivery(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	w := &logger.IOWriter{
		Writer: buf,
	}
	logger.DebugLogger.SetWriter(w)
	d := NewTestDelivery()
	c := notification.NewContent()
	c.Set("contentkey", "contentvalue")
	status, receipt, err := d.Deliver(c)
	if status != notificationdelivery.DeliveryStatusSuccess || err != nil || receipt != "" {
		t.Fatal(status, receipt, err)
	}
	s := buf.String()
	if !strings.Contains(s, "contentkey") || !strings.Contains(s, "contentvalue") {
		t.Fatal(s)
	}
}
