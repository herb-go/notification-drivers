package tencentminiprogramumdelivery

import (
	"testing"

	"github.com/herb-go/notification/notificationdelivery"
)

func NewTestDelivery() *Delivery {
	dc := &notificationdelivery.Config{
		DeliveryType: DeliveryType,
		DeliveryConfig: func(v interface{}) error {
			v.(*Config).App = TestApp
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

func TestTestMessage(t *testing.T) {
	d := NewTestDelivery()
	c := NewTestContent()
	status, receipt, err := d.Deliver(c)
	if status != notificationdelivery.DeliveryStatusSuccess || err != nil || receipt != "" {
		t.Fatal(status, receipt, err)
	}
}
