package tencentcloudsmsdelivery

import (
	"strings"
	"testing"

	"github.com/herb-go/notification/notificationdelivery"
)

func NewTestDelivery() *Delivery {
	dc := &notificationdelivery.Config{
		DeliveryType: DeliveryType,
		DeliveryConfig: func(v interface{}) error {
			v.(*Config).Sms.SdkAppid = TestSMS.SdkAppid
			v.(*Config).Sms.SecretID = TestSMS.SecretID
			v.(*Config).Sms.SecretKey = TestSMS.SecretKey
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
	if status != notificationdelivery.DeliveryStatusSuccess || err != nil || receipt == "" {
		t.Fatal(status, receipt, err)
	}
}

func TestEscape(t *testing.T) {
	msg := "1,2"
	d := NewTestDelivery()
	escaped := d.MustEscape(msg)
	if strings.Contains(escaped, ",") {
		t.Fatal(escaped)
	}
	unescaped := d.Unescape(escaped)
	if unescaped != msg {
		t.Fatal(unescaped, msg)
	}
}
func init() {
}
