package wechatworkdelivery

import "github.com/herb-go/notification/notificationdelivery"

var DeliveryType = "wechatwork"

var ContentNameToUser = "touser"
var ContentNameToTag = "totag"
var ContentNameToParty = "toparty"
var ContentNameMsgType = "msgtype"
var ContentNameContent = "content"
var ContentNameMediaID = "media_id"
var ContentNameMediaFilename = "media_filename"
var ContentNameMediaDataURI = "media_datauri"
var ContentNameTitle = "title"
var ContentNameDescription = "description"
var ContentNameURL = "url"
var ContentNamePicURL = "picurl"
var ContentNameNews = "news"
var ContentNameMPNews = "mpnews"
var ContentNameBtntxt = "btntxt"
var ContentNameTaskID = "task_id"
var ContentNameBtn = "btn"
var ContentNameSafe = "safe"

var RequeiredContent = []string{ContentNameMsgType}

var Fields = []*notificationdelivery.Field{
	{
		Name:    ContentNameToUser,
		Example: "TOUSER",
		Escape:  "",
	},
	{
		Name:    ContentNameToTag,
		Example: "TOTAG",
		Escape:  "",
	},
	{
		Name:    ContentNameURL,
		Example: "https://URL",
		Escape:  "",
	},
	{
		Name:    ContentNameMsgType,
		Example: `MSGTYPE`,
		Escape:  "",
	},
	{
		Name:    ContentNameContent,
		Example: `CONTENT`,
		Escape:  "",
	},
	{
		Name:    ContentNameMediaID,
		Example: `MEDIAID`,
		Escape:  "",
	},
	{
		Name:    ContentNameMediaFilename,
		Example: `FILENAME.EXT`,
		Escape:  "",
	},
	{
		Name:    ContentNameMediaDataURI,
		Example: `https://URL`,
		Escape:  "",
	},
	{
		Name:    ContentNameTitle,
		Example: `TITLE`,
		Escape:  "",
	},
	{
		Name:    ContentNameDescription,
		Example: `DESCRIPTION`,
		Escape:  "",
	},
	{
		Name:    ContentNameURL,
		Example: `https://URL`,
		Escape:  "",
	},
	{
		Name:    ContentNamePicURL,
		Example: `https://URL/1.png`,
		Escape:  "",
	},
	{
		Name:    ContentNameNews,
		Example: `{"articles":[{"title":"github","description":"github Description","url":"https://github.com/","picurl":"https://github.githubassets.com/images/icons/emoji/unicode/1f503.png"},{"title":"only title"}]}`,
		Escape:  "jsonescape",
	},
	{
		Name:    ContentNameMPNews,
		Example: `{"articles":[{"title":"Title","thumb_media_id":"%s","author":"Author","content_source_url":"https://github.com","content":"Content","digest":"Digest description"}]}`,
		Escape:  "jsonescape",
	},
	{
		Name:    ContentNameBtntxt,
		Example: `BTN_TXT`,
		Escape:  "",
	},
	{
		Name:    ContentNameTaskID,
		Example: `TASK_ID`,
		Escape:  "",
	},
	{
		Name:    ContentNameBtn,
		Example: `[{"key":"key111","name":"批准","replace_name":"已批准","color":"red","is_bold":true},{"key":"key222","name":"驳回","replace_name":"已驳回"}]`,
		Escape:  "jsonescape",
	},
	{
		Name:    ContentNameSafe,
		Example: ``,
		Escape:  "",
	},
}
