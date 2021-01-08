package wechattmdelivery

import (
	"testing"

	"github.com/herb-go/notification"
)

func NewTestDelivery() *Delivery {
	d := &Delivery{}
	d.App.AppID = TestAppID
	d.App.AppSecret = TestSecret
	return d
}

var _ notification.Driver = &Delivery{}

func TestTestMessage(t *testing.T) {
	d := NewTestDelivery()
	c := notification.NewContent()
	c.Set(ContentNameToUser, TestTo)
	c.Set(ContentNameTemplateID, TestTemplateID)
	c.Set(ContentNameData, TestData)
	c.Set(ContentNameURL, TestURL)
	status, receipt, err := d.Deliver(c)
	if status != notification.DeliveryStatusSuccess || err != nil || receipt == "" {
		t.Fatal(status, receipt, err)
	}
}

func init() {
}
