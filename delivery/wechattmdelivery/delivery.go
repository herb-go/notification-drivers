package wechattmdelivery

import (
	"encoding/json"
	"net/url"
	"strconv"

	"github.com/herb-go/notification/notificationdelivery"
	"github.com/herb-go/providers/tencent/wechatmp/templatemessage"

	"github.com/herb-go/notification"
	"github.com/herb-go/providers/tencent/wechatmp"
)

type Delivery struct {
	App wechatmp.App
}

func (d *Delivery) buildMsg(c notification.Content) *wechatmp.TemplateMessage {
	msg := wechatmp.NewTemplateMessage()
	msg.ToUser = c.Get(ContentNameToUser)
	msg.TemlpateID = c.Get(ContentNameTemplateID)
	url := c.Get(ContentNameURL)
	if url != "" {
		msg.URL = &url
	}
	miniprogram := c.Get(ContentNameMiniProgram)
	if miniprogram != "" {
		msg.Miniprogram = &wechatmp.TemplateMessageMiniprogram{
			AppID:    ContentNameMiniProgram,
			PagePath: c.Get(ContentNamePagePath),
		}
	}
	msg.Data = json.RawMessage(c.Get(ContentNameData))
	return msg
}
func (d *Delivery) DeliveryType() string {
	return DeliveryType
}
func (d *Delivery) Deliver(c notification.Content) (notificationdelivery.DeliveryStatus, string, error) {
	err := notification.CheckRequiredContentError(c, RequeiredContent)
	if err != nil {
		return notificationdelivery.DeliveryStatusAbort, "", err
	}
	msg := d.buildMsg(c)
	result, err := templatemessage.SendTemplateMessage(&d.App, msg)
	if err != nil {
		return notificationdelivery.DeliveryStatusFail, "", err
	}
	if !result.IsOK() {
		return notificationdelivery.DeliveryStatusFail, "", result
	}
	return notificationdelivery.DeliveryStatusSuccess, strconv.FormatInt(result.MsgID, 10), nil
}

func (d *Delivery) MustEscape(unescaped string) string {
	return url.PathEscape(unescaped)
}

type Config struct {
	wechatmp.App
}

var Factory = func(loader func(interface{}) error) (notificationdelivery.DeliveryDriver, error) {
	c := &Config{}
	err := loader(c)
	if err != nil {
		return nil, err
	}
	d := &Delivery{
		App: c.App,
	}
	return d, nil
}

func init() {
	notificationdelivery.Register(DeliveryType, Factory)
}
