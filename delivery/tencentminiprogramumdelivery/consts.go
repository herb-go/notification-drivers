package tencentminiprogramumdelivery

import "github.com/herb-go/notification/notificationdelivery"

var DeliveryType = "tencentminiprogramum"

var ContentNameToUser = "touser"
var ContentNameAppID = "appid"
var ContentNameTemplateID = "template_id"
var ContentNameURL = "url"
var ContentNameData = "data"
var ContentNamePagePath = "pagepath"
var ContentNameMiniprogram = "miniprogram"
var ContentNameWeappTemplateID = "weapp_template_id"
var ContentNameWeappPage = "weapp_page"
var ContentNameWeappFormID = "weapp_form_id"
var ContentNameWeappEmphasisKeyword = "weapp_emphasis_keyword"
var ContentNameWeappData = "weapp_data"

var RequeiredContent = []string{ContentNameToUser, ContentNameAppID, ContentNameTemplateID}

var Fields = []*notificationdelivery.Field{
	{
		Name:    ContentNameToUser,
		Example: "OPENID",
		Escape:  "",
	},
	{
		Name:    ContentNameAppID,
		Example: "APPID",
		Escape:  "",
	},
	{
		Name:    ContentNameTemplateID,
		Example: "TEMPLATE_ID",
		Escape:  "",
	},
	{
		Name:    ContentNameURL,
		Example: `https://URL`,
		Escape:  "",
	},
	{
		Name:    ContentNameData,
		Example: `{"keyword1":{"value":"339208499","color":"#173177"}}`,
		Escape:  "jsonecape",
	},
	{
		Name:    ContentNameMiniprogram,
		Example: `APPID`,
		Escape:  "",
	},
	{
		Name:    ContentNamePagePath,
		Example: `PAGEPATH`,
		Escape:  "",
	},
	{
		Name:    ContentNameWeappTemplateID,
		Example: `TEMPLATE_ID`,
		Escape:  "",
	},
	{
		Name:    ContentNameWeappPage,
		Example: `PAGE`,
		Escape:  "",
	},
	{
		Name:    ContentNameWeappFormID,
		Example: `FORMID`,
		Escape:  "",
	},
	{
		Name:    ContentNameWeappData,
		Example: `{"keyword1":{"value":"339208499"}}`,
		Escape:  "jsonescape",
	},
	{
		Name:    ContentNameWeappEmphasisKeyword,
		Example: ``,
		Escape:  "",
	},
}
