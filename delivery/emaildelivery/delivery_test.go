package emaildelivery

import (
	"testing"

	"github.com/herb-go/notification"
)

func newAddr(name string, mail string) string {
	d := &Delivery{}
	return d.MustEscape(name) + " <" + d.MustEscape(mail) + ">"
}
func NewTestDelivery() *Delivery {
	d := &Delivery{}
	d.SMTP.Host = TestHost
	d.SMTP.Port = TestPort
	d.SMTP.Identity = TestIdentity
	d.SMTP.Password = TestPassword
	d.SMTP.Username = TestUsername
	d.SMTP.From = TestFrom
	d.SMTP.Sender = TestSender
	return d
}

var _ notification.DeliveryServer = &Delivery{}

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
	status, err := d.Deliver(c)
	if status != notification.DeliveryStatusSuccess || err != nil {
		t.Fatal(status, err)
	}
}
