package wechatworkdelivery

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strconv"
	"testing"
	"time"

	"github.com/herb-go/fetcher"
	"github.com/herb-go/notification/notificationdelivery"

	"github.com/herb-go/notification"
	"github.com/herb-go/providers/tencent/wechatwork"
)

func NewTestDelivery() *Delivery {
	dc := &notificationdelivery.Config{
		DeliveryType: DeliveryType,
		DeliveryConfig: func(v interface{}) error {
			v.(*Config).Agent.AgentID = TestAgentID
			v.(*Config).Agent.CorpID = TestCorpID
			v.(*Config).Agent.Secret = TestSecret
			return nil
		},
	}
	d, err := dc.CreateDriver()
	if err != nil {
		panic(err)
	}
	return d.(*Delivery)
}

var _ notification.DeliveryDriver = &Delivery{}

func TestTextMessage(t *testing.T) {
	d := NewTestDelivery()
	c := notification.NewContent()
	c.Set(ContentNameMsgType, wechatwork.MsgTypeText)
	c.Set(ContentNameToUser, TestRecipient)
	c.Set(ContentNameContent, "test")
	status, receipt, err := d.Deliver(c)
	if status != notification.DeliveryStatusSuccess || err != nil || receipt != "" {
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
	status, receipt, err := d.Deliver(c)
	if status != notification.DeliveryStatusSuccess || err != nil || receipt != "" {
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
	status, receipt, err := d.Deliver(c)
	if status != notification.DeliveryStatusSuccess || err != nil || receipt != "" {
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
	status, receipt, err := d.Deliver(c)
	if status != notification.DeliveryStatusSuccess || err != nil || receipt != "" {
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
	status, receipt, err := d.Deliver(c)
	if status != notification.DeliveryStatusSuccess || err != nil || receipt != "" {
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
	status, receipt, err := d.Deliver(c)
	if status != notification.DeliveryStatusSuccess || err != nil || receipt != "" {
		t.Fatal(status, err)
	}
}

var testMPNews = `
{
	"articles":[
		{
			"title": "Title", 
			"thumb_media_id": "%s",
			"author": "Author",
			"content_source_url": "https://github.com",
			"content": "Content",
			"digest": "Digest description"
		 }
	]
}
`

func TestNewsMPMessage(t *testing.T) {

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

	c.Set(ContentNameMsgType, wechatwork.MsgTypeMPNews)
	c.Set(ContentNameToUser, TestRecipient)
	c.Set(ContentNameMPNews, fmt.Sprintf(testMPNews, mediaid))
	status, receipt, err := d.Deliver(c)
	if status != notification.DeliveryStatusSuccess || err != nil || receipt != "" {
		t.Fatal(status, err)
	}
}

var testBtn = `
[
                {
                    "key": "key111",
                    "name": "批准",
                    "replace_name": "已批准",
                    "color":"red",
                    "is_bold": true
                },
                {
                    "key": "key222",
                    "name": "驳回",
                    "replace_name": "已驳回"
                }
            ]
`

func TestTaskcardMessage(t *testing.T) {
	d := NewTestDelivery()
	c := notification.NewContent()
	c.Set(ContentNameMsgType, wechatwork.MsgTypeTaskcard)
	c.Set(ContentNameToUser, TestRecipient)
	c.Set(ContentNameTitle, "test title")
	c.Set(ContentNameDescription, "test description")
	c.Set(ContentNameURL, "https://github.com/")
	c.Set(ContentNameTaskID, "test-timestamp-"+strconv.FormatInt(time.Now().Unix(), 10))
	c.Set(ContentNameBtn, testBtn)
	status, receipt, err := d.Deliver(c)
	if status != notification.DeliveryStatusSuccess || err != nil || receipt != "" {
		t.Fatal(status, err)
	}
}
func TestTextCardMessage(t *testing.T) {
	d := NewTestDelivery()
	c := notification.NewContent()
	c.Set(ContentNameMsgType, wechatwork.MsgTypeTextcard)
	c.Set(ContentNameToUser, TestRecipient)
	c.Set(ContentNameTitle, "test title")
	c.Set(ContentNameDescription, "test description")
	c.Set(ContentNameBtntxt, "test btn")
	c.Set(ContentNameURL, "https://github.com")
	status, receipt, err := d.Deliver(c)
	if status != notification.DeliveryStatusSuccess || err != nil || receipt != "" {
		t.Fatal(status, err)
	}
}

func TestMarkdownMessage(t *testing.T) {
	d := NewTestDelivery()
	c := notification.NewContent()
	c.Set(ContentNameMsgType, wechatwork.MsgTypeMarkdown)
	c.Set(ContentNameToUser, TestRecipient)
	c.Set(ContentNameContent, `
	# test message

	## sub title
	* value1
	* value2
	`)
	status, receipt, err := d.Deliver(c)
	if status != notification.DeliveryStatusSuccess || err != nil || receipt != "" {
		t.Fatal(status, err)
	}
}
func init() {
}
