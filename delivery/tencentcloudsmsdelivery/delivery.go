package tencentcloudsmsdelivery

import (
	"html"
	"strings"

	"github.com/herb-go/notification/notificationdelivery"
	"github.com/herb-go/providers/tencent/tencentcloud/tencentcloudsms"

	"github.com/herb-go/notification"
)

type Delivery struct {
	Sms tencentcloudsms.Sms
}

//CheckInvalidContent check if given content invalid
//Return invalid fields and any error raised
func (d *Delivery) CheckInvalidContent(c notification.Content) ([]string, error) {
	invalids := notification.CheckRequiredContent(c, RequeiredContent)
	if len(invalids) > 0 {
		return invalids, nil
	}
	return []string{}, nil
}

func (d *Delivery) DeliveryType() string {
	return DeliveryType
}
func (d *Delivery) buildMsg(c notification.Content) (*tencentcloudsms.Message, error) {
	m := tencentcloudsms.NewMessage()
	m.TemplateID = c.Get(ContentNameTemplateID)
	pn := c.Get(ContentNamePhoneNumber)
	if pn != "" {
		list := strings.Split(pn, ",")
		m.PhoneNumber = make([]string, len(list))
		for k := range list {
			m.PhoneNumber[k] = d.Unescape(list[k])
		}
	}
	tp := c.Get(ContentNameTemplateParam)
	if tp != "" {
		list := strings.Split(tp, ",")
		m.TemplateParam = make([]string, len(list))
		for k := range list {
			m.TemplateParam[k] = d.Unescape(list[k])
		}
	}
	m.Sign = c.Get(ContentNameSign)
	m.SessionContext = c.Get(ContentNameSessionContext)
	m.ExtendCode = c.Get(ContentNameExtendCode)
	m.SenderID = c.Get(ContentNameSenderID)
	return m, nil
}

//ContentFields return content fields
//Return invalid fields and any error raised
func (d *Delivery) ContentFields() []*notificationdelivery.Field {
	return Fields
}

func (d *Delivery) Deliver(c notification.Content) (notificationdelivery.DeliveryStatus, string, error) {
	err := notification.CheckRequiredContentError(c, RequeiredContent)
	if err != nil {
		return notificationdelivery.DeliveryStatusAbort, "", err
	}
	msg, err := d.buildMsg(c)
	if err != nil {
		return notificationdelivery.DeliveryStatusAbort, "", err
	}
	result, err := d.Sms.Send(msg)
	if err != nil {
		return notificationdelivery.DeliveryStatusFail, "", err
	}
	return notificationdelivery.DeliveryStatusSuccess, result.Response.RequestId, nil
}
func (d *Delivery) Unescape(escaped string) string {
	return html.UnescapeString(escaped)
}

var escaper = strings.NewReplacer(
	`,`, "&#44;",
	`&`, "&amp;",
)

func (d *Delivery) MustEscape(unescaped string) string {
	return escaper.Replace(unescaped)
}

type Config struct {
	tencentcloudsms.Sms
}

var Factory = func(loader func(interface{}) error) (notificationdelivery.DeliveryDriver, error) {
	c := &Config{}
	err := loader(c)
	if err != nil {
		return nil, err
	}
	d := &Delivery{
		Sms: c.Sms,
	}
	return d, nil
}

func init() {
	notificationdelivery.Register(DeliveryType, Factory)
}
