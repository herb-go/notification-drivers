package wechattmdelivery

import (
	"testing"

	"github.com/herb-go/notification"
	"github.com/herb-go/notification/notificationdelivery"
)

func NewTestDelivery() *Delivery {
	dc := &notificationdelivery.Config{
		DeliveryType: DeliveryType,
		DeliveryConfig: func(v interface{}) error {
			c := &Config{}
			c.App.AppID = TestAppID
			c.App.AppSecret = TestSecret

			v.(*Config).App = c.App
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
	c := notification.NewContent()
	c.Set(ContentNameToUser, TestTo)
	c.Set(ContentNameTemplateID, TestTemplateID)
	c.Set(ContentNameData, TestData)
	c.Set(ContentNameURL, TestURL)
	status, receipt, err := d.Deliver(c)
	if status != notificationdelivery.DeliveryStatusSuccess || err != nil || receipt == "" {
		t.Fatal(status, receipt, err)
	}
}

func init() {
}
