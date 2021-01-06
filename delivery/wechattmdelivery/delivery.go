package wechattmdelivery

import (
	"encoding/json"
	"net/url"

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
func (d *Delivery) Deliver(c notification.Content) (notification.DeliveryStatus, error) {
	err := notification.CheckRequiredContentError(c, RequeiredContent)
	if err != nil {
		return notification.DeliveryStatusAbort, err
	}
	msg := d.buildMsg(c)
	result, err := templatemessage.SendTemplateMessage(&d.App, msg)
	if err != nil {
		return notification.DeliveryStatusFail, err
	}
	if !result.IsOK() {
		return notification.DeliveryStatusFail, result
	}
	return notification.DeliveryStatusSuccess, nil
}

func (d *Delivery) MustEscape(unescaped string) string {
	return url.PathEscape(unescaped)
}
