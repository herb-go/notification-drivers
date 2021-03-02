package emaildelivery

import "github.com/herb-go/notification/notificationdelivery"

var DeliveryType = "email"

var ContentNameSender = "sender"
var ContentNameFrom = "from"
var ContentNameSubject = "subject"
var ContentNameText = "text"
var ContentNameHTML = "html"
var ContentNameReplyTo = "replyto"
var ContentNameTo = "to"
var ContentNameCC = "cc"
var ContentNameBCC = "bcc"
var ContentNameAttachments = "attachments"
var RequeiredContent = []string{ContentNameTo}

var Fields = []*notificationdelivery.Field{
	{
		Name:    ContentNameFrom,
		Example: "NAME <mail@example.com>",
		Escape:  "",
	},
	{
		Name:    ContentNameSubject,
		Example: "EMAIL SUBJECT",
		Escape:  "",
	},
	{
		Name:    ContentNameTo,
		Example: `TONAME <to@example.com>,mail@example.com`,
		Escape:  "commaescape",
	},
	{
		Name:    ContentNameText,
		Example: "THIS IS A TEXT EMAIL\nLINE 2.",
		Escape:  "",
	},
	{
		Name:    ContentNameHTML,
		Example: "<p>\nTHIS IS A <b>HTML</b> EMAIL\n</p>",
		Escape:  "",
	},
	{
		Name:    ContentNameAttachments,
		Example: "[{\"Filename\":\"1.png\",\"DataURI\":\"https://URL\",\"ContentType\":\"image/png\"},{\"Filename\":\"2.png\",\"DataURI\":\"https://URL2\"}]",
		Escape:  "jsonescape",
	},
	{
		Name:    ContentNameReplyTo,
		Example: `REPLYTONAME <replyto@example.com>`,
		Escape:  "",
	},
	{
		Name:    ContentNameCC,
		Example: `CCNAME <cc@example.com>,mail@example.com`,
		Escape:  "commaescape",
	},
	{
		Name:    ContentNameBCC,
		Example: `BCCNAME <bcc@example.com>,mail@example.com`,
		Escape:  "commaescape",
	},
	{
		Name:    ContentNameSender,
		Example: "",
		Escape:  "",
	},
}
