package aliyunsmsdelivery

import (
	"html"
	"strings"

	"github.com/herb-go/notification/notificationdelivery"
	"github.com/herb-go/providers/alibaba/aliyun"
	"github.com/herb-go/providers/alibaba/aliyun/aliyunsms"

	"github.com/herb-go/notification"
)

type Delivery struct {
	*aliyun.AccessKey
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
func (d *Delivery) buildMsg(c notification.Content) (*aliyunsms.Message, error) {
	m := aliyunsms.NewMessage()
	m.PhoneNumbers = c.Get(ContentNamePhoneNumbers)
	m.SignName = c.Get(ContentNameSignName)
	m.TemplateCode = c.Get(ContentNameTemplateCode)
	m.OutID = c.Get(ContentNameOutID)
	m.TemplateParam = c.Get(ContentNameTemplateParam)
	m.SmsUpExtendCode = c.Get(ContentNameSmsUpExtendCode)
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
	result, err := aliyunsms.Send(d.AccessKey, msg)
	if err != nil {
		return notificationdelivery.DeliveryStatusFail, "", err
	}
	return notificationdelivery.DeliveryStatusSuccess, result.BizId, nil
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
	*aliyun.AccessKey
}

var Factory = func(loader func(interface{}) error) (notificationdelivery.DeliveryDriver, error) {
	c := &Config{
		AccessKey: &aliyun.AccessKey{},
	}
	err := loader(c)
	if err != nil {
		return nil, err
	}
	d := &Delivery{
		AccessKey: c.AccessKey,
	}
	return d, nil
}

func init() {
	notificationdelivery.Register(DeliveryType, Factory)
}
