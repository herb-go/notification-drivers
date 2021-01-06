package wechatworkdelivery

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/herb-go/fetcher"

	"github.com/herb-go/notification"
	"github.com/herb-go/providers/tencent/wechatwork"
)

func NewTestDelivery() *Delivery {
	d := &Delivery{}
	d.Agent.AgentID = TestAgentID
	d.Agent.CorpID = TestCorpID
	d.Agent.Secret = TestSecret
	return d
}

var _ notification.DeliveryServer = &Delivery{}

func TestTestMessage(t *testing.T) {
	d := NewTestDelivery()
	c := notification.NewContent()
	c.Set(ContentNameMsgType, wechatwork.MsgTypeText)
	c.Set(ContentNameToUser, TestRecipient)
	c.Set(ContentNameContent, "test")
	status, err := d.Deliver(c)
	if status != notification.DeliveryStatusSuccess || err != nil {
		t.Fatal(status, err)
	}
}

func TestImageMessage(t *testing.T) {
	d := NewTestDelivery()
	c := notification.NewContent()
	var result []byte
	_, err := fetcher.NewPreset().With(fetcher.URL(TestPictureURL)).FetchAndParse(fetcher.Should200(fetcher.AsBytes(&result)))
	if err != nil {
		panic(err)
	}
	mediaid, err := d.Agent.MediaUpload(wechatwork.MediaTypeImage, "test.png", bytes.NewBuffer(result))
	if err != nil {
		panic(err)
	}
	c.Set(ContentNameMsgType, wechatwork.MsgTypeImage)
	c.Set(ContentNameToUser, TestRecipient)
	c.Set(ContentNameMediaID, mediaid)
	status, err := d.Deliver(c)
	if status != notification.DeliveryStatusSuccess || err != nil {
		t.Fatal(status, err)
	}
}

func TestVoiceMessage(t *testing.T) {
	d := NewTestDelivery()
	c := notification.NewContent()
	if TestAmrFile == "" {
		t.Fatal()
	}
	result, err := ioutil.ReadFile(TestAmrFile)
	if err != nil {
		panic(err)
	}
	mediaid, err := d.Agent.MediaUpload(wechatwork.MediaTypeVoice, "test.mp3", bytes.NewBuffer(result))
	if err != nil {
		panic(err)
	}
	c.Set(ContentNameMsgType, wechatwork.MsgTypeVoice)
	c.Set(ContentNameToUser, TestRecipient)
	c.Set(ContentNameMediaID, mediaid)
	status, err := d.Deliver(c)
	if status != notification.DeliveryStatusSuccess || err != nil {
		t.Fatal(status, err)
	}
}

func TestVideoMessage(t *testing.T) {
	d := NewTestDelivery()
	c := notification.NewContent()
	if TestAmrFile == "" {
		t.Fatal()
	}
	result, err := ioutil.ReadFile(TestMp4File)
	if err != nil {
		panic(err)
	}
	mediaid, err := d.Agent.MediaUpload(wechatwork.MediaTypeVideo, "test.mp4", bytes.NewBuffer(result))
	if err != nil {
		panic(err)
	}
	c.Set(ContentNameMsgType, wechatwork.MsgTypeVideo)
	c.Set(ContentNameToUser, TestRecipient)
	c.Set(ContentNameMediaID, mediaid)
	c.Set(ContentNameTitle, "video title")
	c.Set(ContentNameDescription, "video desc")
	status, err := d.Deliver(c)
	if status != notification.DeliveryStatusSuccess || err != nil {
		t.Fatal(status, err)
	}
}

func TestFileMessage(t *testing.T) {
	d := NewTestDelivery()
	c := notification.NewContent()
	if TestAmrFile == "" {
		t.Fatal()
	}
	result, err := ioutil.ReadFile(TestMp4File)
	if err != nil {
		panic(err)
	}
	mediaid, err := d.Agent.MediaUpload(wechatwork.MediaTypeFile, "test.mp4", bytes.NewBuffer(result))
	if err != nil {
		panic(err)
	}
	c.Set(ContentNameMsgType, wechatwork.MsgTypeFile)
	c.Set(ContentNameToUser, TestRecipient)
	c.Set(ContentNameMediaID, mediaid)
	status, err := d.Deliver(c)
	if status != notification.DeliveryStatusSuccess || err != nil {
		t.Fatal(status, err)
	}
}

var testNews = `
{
"articles":[
	{
		"title": "github",
		"description": "github Description",
		"url": "https://github.com/",
		"picurl": "https://github.githubassets.com/images/icons/emoji/unicode/1f503.png"
	},
	{
		"title": "only title"
	}    
]
	}
`

func TestNewsMessage(t *testing.T) {
	d := NewTestDelivery()
	c := notification.NewContent()
	c.Set(ContentNameMsgType, wechatwork.MsgTypeNews)
	c.Set(ContentNameToUser, TestRecipient)
	c.Set(ContentNameNews, testNews)
	status, err := d.Deliver(c)
	if status != notification.DeliveryStatusSuccess || err != nil {
		t.Fatal(status, err)
	}
}
func init() {
}
