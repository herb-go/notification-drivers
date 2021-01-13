package emaildelivery

import (
	"mime"
	"net/smtp"
	"strconv"
	"strings"

	"github.com/herb-go/notification"
	"github.com/herb-go/notification/notificationdelivery"

	"github.com/jordan-wright/email"
)

type SMTP struct {
	Sender string
	// Host smtp host addr
	Host string

	Port int
	// Identity user identity(user account) for stmp arddr
	Identity string
	// email from address
	From string
	// Username email stmp user name
	Username string
	// Pasword email smtp password
	Password string
	StartTLS bool
}

func (s *SMTP) NewEmail(c notification.Content) *email.Email {
	msg := email.NewEmail()
	msg.From = s.From
	if msg.From == "" {
		msg.From = c.Get(ContentNameFrom)
	}
	msg.Sender = s.Sender
	if msg.Sender == "" {
		msg.Sender = c.Get(ContentNameSender)
	}
	msg.Subject = c.Get(ContentNameSubject)
	text := c.Get(ContentNameText)
	if text != "" {
		msg.Text = []byte(text)
	}
	html := c.Get(ContentNameHTML)
	if html != "" {
		msg.HTML = []byte(html)
	}
	replyto := c.Get(ContentNameReplyTo)
	if replyto != "" {
		msg.ReplyTo = strings.Split(replyto, Separator)
	}
	to := c.Get(ContentNameTo)
	if to != "" {
		msg.To = strings.Split(to, Separator)
	}
	cc := c.Get(ContentNameCC)
	if cc != "" {
		msg.Cc = strings.Split(cc, Separator)
	}
	bcc := c.Get(ContentNameBCC)
	if bcc != "" {
		msg.Bcc = strings.Split(bcc, Separator)
	}
	return msg
}

func (s *SMTP) Send(msg *email.Email) error {
	if s.StartTLS {
		msg.SendWithStartTLS(s.Host+":"+strconv.Itoa(s.Port), smtp.PlainAuth(s.Identity, s.Username, s.Password, s.Host), nil)
	}
	return msg.Send(s.Host+":"+strconv.Itoa(s.Port), smtp.PlainAuth(s.Identity, s.Username, s.Password, s.Host))
}

type Delivery struct {
	SMTP SMTP
}

func (d *Delivery) DeliveryType() string {
	return DeliveryType
}
func (d *Delivery) Deliver(c notification.Content) (notification.DeliveryStatus, string, error) {
	err := notification.CheckRequiredContentError(c, RequeiredContent)
	if err != nil {
		return notification.DeliveryStatusAbort, "", err
	}
	msg := d.SMTP.NewEmail(c)
	err = d.SMTP.Send(msg)
	if err != nil {
		return notification.DeliveryStatusFail, "", err
	}
	return notification.DeliveryStatusSuccess, "", nil

}

func (d *Delivery) MustEscape(unescaped string) string {
	return mime.BEncoding.Encode("utf-8", unescaped)
}

type Config struct {
	SMTP SMTP
}

var Factory = func(loader func(interface{}) error) (notification.DeliveryDriver, error) {
	c := &Config{}
	err := loader(c)
	if err != nil {
		return nil, err
	}
	d := &Delivery{
		SMTP: c.SMTP,
	}
	return d, nil
}

func init() {
	notificationdelivery.Register(DeliveryType, Factory)
}
