package emaildelivery

import (
	"testing"

	"github.com/herb-go/notification"
	"github.com/herb-go/notification/notificationdelivery"
)

func newAddr(name string, mail string) string {
	d := &Delivery{}
	return d.MustEscape(name) + " <" + d.MustEscape(mail) + ">"
}
func NewTestDelivery() *Delivery {
	c := &Config{}
	c.SMTP.Host = TestHost
	c.SMTP.Port = TestPort
	c.SMTP.Identity = TestIdentity
	c.SMTP.Password = TestPassword
	c.SMTP.Username = TestUsername
	c.SMTP.From = TestFrom
	c.SMTP.Sender = TestSender
	c.SMTP.StartTLS = TestStartTLS
	dc := &notificationdelivery.Config{
		DeliveryType: DeliveryType,
		DeliveryConfig: func(v interface{}) error {
			v.(*Config).SMTP = c.SMTP
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
	d := NewTestDelivery()
	c := notification.NewContent()
	c.Set(ContentNameTo, TestTo)
	c.Set(ContentNameReplyTo, TestReplyTO)
	c.Set(ContentNameCC, TestCC)
	c.Set(ContentNameBCC, TestBCC)
	c.Set(ContentNameSubject, " test subject 😅")
	c.Set(ContentNameText, "text body")
	c.Set(ContentNameHTML, "<p><b>html</b> body</p>")
	status, receipt, err := d.Deliver(c)
	if status != notificationdelivery.DeliveryStatusSuccess || err != nil || receipt != "" {
		t.Fatal(status, receipt, err)
	}
}
