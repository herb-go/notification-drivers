package wechattmdelivery

import "github.com/herb-go/notification/notificationdelivery"

var DeliveryType = "wechattm"

var ContentNameToUser = "touser"
var ContentNameTemplateID = "template_id"
var ContentNameURL = "url"
var ContentNameMiniProgram = "miniprogram"
var ContentNamePagePath = "pagepath"
var ContentNameData = "data"

var RequeiredContent = []string{ContentNameToUser, ContentNameTemplateID, ContentNameData}

var Fields = []*notificationdelivery.Field{
	{
		Name:    ContentNameToUser,
		Example: "TOUSER",
		Escape:  "",
	},
	{
		Name:    ContentNameTemplateID,
		Example: "TEMPLATE_ID",
		Escape:  "",
	},
	{
		Name:    ContentNameURL,
		Example: "https://URL",
		Escape:  "",
	},
	{
		Name:    ContentNameMiniProgram,
		Example: `APPID`,
		Escape:  "",
	},
	{
		Name:    ContentNamePagePath,
		Example: `PAGE`,
		Escape:  "",
	},
	{
		Name:    ContentNameData,
		Example: `{\"first\": {\"value\":\"VALUE\",\"color\":\"#173177\"}}`,
		Escape:  "jsonecape",
	},
}
