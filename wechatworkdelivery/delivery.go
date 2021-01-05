package wechatworkdelivery

import (
	"encoding/json"

	"github.com/herb-go/notification"
	"github.com/herb-go/providers/tencent/wechatwork"
)

type Delivery struct {
	ID    string
	Agent wechatwork.Agent
}

func (d *Delivery) DeliveryID() string {
	return d.ID
}
func (d *Delivery) DeliveryType() string {
	return DeliveryType
}
func (d *Delivery) Deliver(c notification.Content) (notification.DeliveryStatus, error) {
	err := notification.CheckRequiredContentError(c, RequeiredContent)
	if err != nil {
		return notification.DeliveryStatusAbort, err
	}
	var msgtype = c.Get(ContentNameMsgType)
	if msgtype == "" {
		return notification.DeliveryStatusAbort, NewInvalidMsgType(msgtype)
	}
	msg := wechatwork.NewMessage()
	var initFunc func(msg *wechatwork.Message, c notification.Content) error
	switch msgtype {
	case wechatwork.MsgTypeText:
		initFunc = d.initTextMessage
	case wechatwork.MsgTypeImage:
		initFunc = d.initImageMessage
	case wechatwork.MsgTypeVoice:
		initFunc = d.initVoiceMessage
	case wechatwork.MsgTypeVideo:
		initFunc = d.initVideoMessage
	case wechatwork.MsgTypeFile:
		initFunc = d.initFileMessage
	case wechatwork.MsgTypeNews:
		initFunc = d.initNewsMessage

	default:
		return notification.DeliveryStatusAbort, NewInvalidMsgType(msgtype)
	}
	err = initFunc(msg, c)
	if err != nil {
		return notification.DeliveryStatusFail, err
	}
	msg.MsgType = msgtype
	d.initMessage(msg, c)
	result, err := d.Agent.SendMessage(msg)
	if err != nil {
		return notification.DeliveryStatusFail, err
	}
	if result.Errcode != 0 {
		return notification.DeliveryStatusAbort, wechatwork.NewResultError(result.Errcode, result.Errmsg)
	}
	return notification.DeliveryStatusSuccess, nil

}

func (d *Delivery) initMessage(msg *wechatwork.Message, c notification.Content) {
	touser := c.Get(ContentNameToUser)
	if touser != "" {
		msg.ToUser = &touser
	}
	toparty := c.Get(ContentNameToParty)
	if toparty != "" {
		msg.ToParty = &toparty
	}
	totag := c.Get(ContentNameToTag)
	if totag != "" {
		msg.ToTag = &totag
	}
	safe := c.Get(ContentNameSafe)
	if safe != "" {
		msg.Safe = 1
	}
	msg.AgentID = d.Agent.AgentID
}
func (d *Delivery) initTextMessage(msg *wechatwork.Message, c notification.Content) error {
	msg.Text = &wechatwork.MessageText{
		Content: c.Get(ContentNameContent),
	}
	return nil
}
func (d *Delivery) initImageMessage(msg *wechatwork.Message, c notification.Content) error {
	msg.Image = &wechatwork.MessageMedia{
		MediaID: c.Get(ContentNameMediaID),
	}
	return nil
}
func (d *Delivery) initVoiceMessage(msg *wechatwork.Message, c notification.Content) error {
	msg.Voice = &wechatwork.MessageMedia{
		MediaID: c.Get(ContentNameMediaID),
	}
	return nil
}
func (d *Delivery) initVideoMessage(msg *wechatwork.Message, c notification.Content) error {
	msg.Video = &wechatwork.MessageVideo{
		MediaID: c.Get(ContentNameMediaID),
	}
	title := c.Get(ContentNameTitle)
	if title != "" {
		msg.Video.Title = &title
	}
	desc := c.Get(ContentNameDescription)
	if desc != "" {
		msg.Video.Description = &desc
	}
	return nil
}

func (d *Delivery) initNewsMessage(msg *wechatwork.Message, c notification.Content) error {
	var news = &wechatwork.MessageNews{}
	err := json.Unmarshal([]byte(c.Get(ContentNameNews)), news)
	if err != nil {
		return ErrNewsFormatError
	}
	msg.News = news
	return nil
}

func (d *Delivery) initFileMessage(msg *wechatwork.Message, c notification.Content) error {
	msg.File = &wechatwork.MessageMedia{
		MediaID: c.Get(ContentNameMediaID),
	}
	return nil
}
