package emaildelivery

import (
	"bytes"
	"encoding/json"
	"mime"
	"net/mail"
	"net/smtp"
	"strconv"
	"strings"

	"github.com/herb-go/herbtext-drivers/commonenvironment"

	"github.com/herb-go/herbdata/datauri"
	"github.com/herb-go/notification"
	"github.com/herb-go/notification/notificationdelivery"

	"github.com/jordan-wright/email"
)

type Attachment struct {
	Filename    string
	DataURI     string
	ContentType string
}

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

func converntMailList(addrs string) ([]string, error) {

	addrlist := strings.Split(addrs, ",")
	result := make([]string, len(addrlist))
	for k, v := range addrlist {
		if v == "" {
			continue
		}
		addr, err := mail.ParseAddress(commonenvironment.ConverterCommaUnescape(v))
		if err != nil {
			return nil, err
		}
		result[k] = addr.String()
	}
	return result, nil
}
func (s *SMTP) NewEmail(c notification.Content) (*email.Email, error) {
	var err error
	msg := email.NewEmail()
	from := c.Get(ContentNameFrom)
	fromlist, err := converntMailList(from)
	if err != nil {
		return nil, err
	}
	msg.From = fromlist[0]
	if msg.From == "" {
		msg.From = s.From
	}
	msg.Sender = c.Get(ContentNameSender)
	if msg.Sender == "" {
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
		msg.ReplyTo, err = converntMailList(replyto)
		if err != nil {
			return nil, err
		}
	}
	to := c.Get(ContentNameTo)
	if to != "" {
		msg.To, err = converntMailList(to)
		if err != nil {
			return nil, err
		}

	}
	cc := c.Get(ContentNameCC)
	if cc != "" {
		msg.Cc, err = converntMailList(cc)
		if err != nil {
			return nil, err
		}
	}
	bcc := c.Get(ContentNameBCC)
	if bcc != "" {
		msg.Bcc, err = converntMailList(bcc)
		if err != nil {
			return nil, err
		}
	}
	attachmentsjson := c.Get(ContentNameAttachments)
	if attachmentsjson != "" {
		attachmentlist := []*Attachment{}
		err := json.Unmarshal([]byte(attachmentsjson), &attachmentlist)
		if err != nil {
			return nil, err
		}
		for _, v := range attachmentlist {
			data, err := datauri.Load(v.DataURI)
			if err != nil {
				return nil, err
			}
			_, err = msg.Attach(bytes.NewBuffer(data), v.Filename, v.ContentType)
			if err != nil {
				return nil, err
			}
		}
	}
	return msg, nil
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
func (d *Delivery) Deliver(c notification.Content) (notificationdelivery.DeliveryStatus, string, error) {
	_, err := d.CheckInvalidContent(c)
	if err != nil {
		return notificationdelivery.DeliveryStatusAbort, "", err
	}
	msg, err := d.SMTP.NewEmail(c)
	if err != nil {
		return notificationdelivery.DeliveryStatusFail, "", err
	}
	err = d.SMTP.Send(msg)
	if err != nil {
		return notificationdelivery.DeliveryStatusFail, "", err
	}
	return notificationdelivery.DeliveryStatusSuccess, "", nil

}

func (d *Delivery) MustEscape(unescaped string) string {
	return mime.BEncoding.Encode("utf-8", unescaped)
}

//ContentFields return content fields
//Return invalid fields and any error raised
func (d *Delivery) ContentFields() []*notificationdelivery.Field {
	return Fields
}

//CheckInvalidContent check if given content invalid
//Return invalid fields and any error raised
func (d *Delivery) CheckInvalidContent(c notification.Content) ([]string, error) {
	result := notification.CheckRequiredContent(c, RequeiredContent)
	if len(result) > 0 {
		return result, nil
	}
	if d.SMTP.From == "" && c.Get(ContentNameFrom) == "" {
		return []string{"from"}, nil
	}
	return nil, nil
}

type Config struct {
	SMTP SMTP
}

var Factory = func(loader func(interface{}) error) (notificationdelivery.DeliveryDriver, error) {
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
