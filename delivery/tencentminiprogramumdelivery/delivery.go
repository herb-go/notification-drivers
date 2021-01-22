package tencentminiprogramumdelivery

import (
	"encoding/json"
	"net/url"

	"github.com/herb-go/notification/notificationdelivery"
	"github.com/herb-go/providers/tencent/tencentminiprogram"
	"github.com/herb-go/providers/tencent/tencentminiprogram/tencentminiprogramum"

	"github.com/herb-go/notification"
)

type Delivery struct {
	App *tencentminiprogram.App
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
func (d *Delivery) buildMsg(c notification.Content) (*tencentminiprogramum.Message, error) {
	m := tencentminiprogramum.NewMessage()
	m.ToUser = c.Get(ContentNameToUser)
	m.MpTemplateMsg.AppID = c.Get(ContentNameAppID)
	m.MpTemplateMsg.TemplateID = c.Get(ContentNameTemplateID)
	url := c.Get(ContentNameURL)
	if url != "" {
		m.MpTemplateMsg.URL = &url
	}
	miniprogram := c.Get(ContentNameMiniprogram)
	if miniprogram != "" {
		m.MpTemplateMsg.Miniprogram = &tencentminiprogramum.TemplateMessageMiniprogram{
			AppID:    miniprogram,
			PagePath: c.Get(ContentNamePagePath),
		}
	}
	m.MpTemplateMsg.Data = json.RawMessage(c.Get(ContentNameData))
	weappTemplateID := c.Get(ContentNameWeappTemplateID)
	if weappTemplateID != "" {
		m.WeappTemplateMessage = &tencentminiprogramum.WeappTemplateMessage{
			TemplateID:      weappTemplateID,
			Page:            c.Get(ContentNameWeappPage),
			FormID:          c.Get(ContentNameWeappFormID),
			EmphasisKeyword: c.Get(ContentNameWeappEmphasisKeyword),
			Data:            c.Get(ContentNameWeappData),
		}
	}
	return m, nil
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
	err = tencentminiprogramum.Send(d.App, msg)
	if err != nil {
		return notificationdelivery.DeliveryStatusFail, "", err
	}
	return notificationdelivery.DeliveryStatusSuccess, "", nil
}

func (d *Delivery) MustEscape(unescaped string) string {
	return url.PathEscape(unescaped)
}

type Config struct {
	*tencentminiprogram.App
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
