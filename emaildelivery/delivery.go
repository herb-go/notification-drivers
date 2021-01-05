package emaildelivery

import (
	"net/smtp"
	"strconv"
	"strings"

	"github.com/herb-go/notification"

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
}

func (s *SMTP) NewEmail(c notification.Content) *email.Email {
	msg := email.NewEmail()
	from := c.Get(ContentNameFrom)
	if from != "" {
		msg.From = from
	} else {
		msg.From = s.From
	}
	sender := c.Get(ContentNameSender)
	if sender != "" {
		msg.Sender = sender
	} else {
		msg.Sender = s.Sender
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
	return msg.Send(s.Host+":"+strconv.Itoa(s.Port), smtp.PlainAuth(s.Identity, s.Username, s.Password, s.Host))
}

type Delivery struct {
	SMTP SMTP
}

func (d *Delivery) DeliveryType() string {
	return DeliveryType
}
func (d *Delivery) Deliver(c notification.Content) (notification.DeliveryStatus, error) {
	err := notification.CheckRequiredContentError(c, RequeiredContent)
	if err != nil {
		return notification.DeliveryStatusAbort, err
	}
	msg := d.SMTP.NewEmail(c)
	err = d.SMTP.Send(msg)
	if err != nil {
		return notification.DeliveryStatusFail, err
	}
	return notification.DeliveryStatusSuccess, nil

}
