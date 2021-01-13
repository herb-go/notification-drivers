package wechatworkdelivery

import (
	"encoding/json"
	"net/url"

	"github.com/herb-go/notification"
	"github.com/herb-go/notification/notificationdelivery"
	"github.com/herb-go/providers/tencent/wechatwork"
)

type Delivery struct {
	Agent wechatwork.Agent
}

func (d *Delivery) DeliveryType() string {
	return DeliveryType
}
func (d *Delivery) Deliver(c notification.Content) (notificationdelivery.DeliveryStatus, string, error) {
	err := notification.CheckRequiredContentError(c, RequeiredContent)
	if err != nil {
		return notificationdelivery.DeliveryStatusAbort, "", err
	}
	var msgtype = c.Get(ContentNameMsgType)
	if msgtype == "" {
		return notificationdelivery.DeliveryStatusAbort, "", NewInvalidMsgType(msgtype)
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
	case wechatwork.MsgTypeMPNews:
		initFunc = d.initMPNewsMessage
	case wechatwork.MsgTypeTextcard:
		initFunc = d.initTextcardMessage
	case wechatwork.MsgTypeMarkdown:
		initFunc = d.initMarkdownMessage
	case wechatwork.MsgTypeTaskcard:
		initFunc = d.initTaskcardMessage

	default:
		return notificationdelivery.DeliveryStatusAbort, "", NewInvalidMsgType(msgtype)
	}
	err = initFunc(msg, c)
	if err != nil {
		return notificationdelivery.DeliveryStatusFail, "", err
	}
	msg.MsgType = msgtype
	d.initMessage(msg, c)
	result, err := d.Agent.SendMessage(msg)
	if err != nil {
		return notificationdelivery.DeliveryStatusFail, "", err
	}
	if result.Errcode != 0 {
		return notificationdelivery.DeliveryStatusAbort, "", wechatwork.NewResultError(result.Errcode, result.Errmsg)
	}
	return notificationdelivery.DeliveryStatusSuccess, "", nil

}
func (d *Delivery) MustEscape(unescaped string) string {
	return url.PathEscape(unescaped)
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
func (d *Delivery) initMPNewsMessage(msg *wechatwork.Message, c notification.Content) error {
	var mpnews = &wechatwork.MessageMPNews{}
	err := json.Unmarshal([]byte(c.Get(ContentNameMPNews)), mpnews)
	if err != nil {
		return ErrMPNewsFormatError
	}
	msg.MPNews = mpnews
	return nil
}
func (d *Delivery) initFileMessage(msg *wechatwork.Message, c notification.Content) error {
	msg.File = &wechatwork.MessageMedia{
		MediaID: c.Get(ContentNameMediaID),
	}
	return nil
}

func (d *Delivery) initTextcardMessage(msg *wechatwork.Message, c notification.Content) error {
	var textcard = &wechatwork.MessageTextcard{}
	textcard.Title = c.Get(ContentNameTitle)
	textcard.Description = c.Get(ContentNameDescription)
	textcard.URL = c.Get(ContentNameURL)
	btntxt := c.Get(ContentNameBtntxt)
	if btntxt != "" {
		textcard.Btntxt = &btntxt
	}
	msg.Textcard = textcard
	return nil
}

func (d *Delivery) initTaskcardMessage(msg *wechatwork.Message, c notification.Content) error {
	var taskcard = &wechatwork.MessageTaskcard{}
	taskcard.Title = c.Get(ContentNameTitle)
	taskcard.Description = c.Get(ContentNameDescription)
	url := c.Get(ContentNameURL)
	taskcard.URL = &url
	taskcard.TaskID = c.Get(ContentNameTaskID)
	btnjson := c.Get(ContentNameBtn)
	var btn []*wechatwork.MessageTaskcardBtn
	err := json.Unmarshal([]byte(btnjson), &btn)
	if err != nil {
		return ErrTaskcardBtnFormatError
	}
	taskcard.Btn = btn
	msg.Taskcard = taskcard
	return nil
}

func (d *Delivery) initMarkdownMessage(msg *wechatwork.Message, c notification.Content) error {
	msg.Markdown = &wechatwork.MessageMarkdown{
		Content: c.Get(ContentNameContent),
	}
	return nil
}

type Config struct {
	wechatwork.Agent
}

var Factory = func(loader func(interface{}) error) (notificationdelivery.DeliveryDriver, error) {
	c := &Config{}
	err := loader(c)
	if err != nil {
		return nil, err
	}
	d := &Delivery{
		Agent: c.Agent,
	}
	return d, nil
}

func init() {
	notificationdelivery.Register(DeliveryType, Factory)
}
