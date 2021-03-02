package aliyunsmsdelivery

import "github.com/herb-go/notification/notificationdelivery"

var DeliveryType = "aliyunsms"

var ContentNamePhoneNumbers = "phonenumbers"
var ContentNameSignName = "signname"
var ContentNameTemplateCode = "templatecode"
var ContentNameTemplateParam = "templateparam"
var ContentNameSmsUpExtendCode = "smsupextendcode"
var ContentNameOutID = "outid"

var RequeiredContent = []string{ContentNamePhoneNumbers, ContentNameTemplateCode}

var Fields = []*notificationdelivery.Field{
	{
		Name:    ContentNamePhoneNumbers,
		Example: "13800000000",
		Escape:  "",
	},
	{
		Name:    ContentNameSignName,
		Example: "COMPANYSIGN",
		Escape:  "",
	},
	{
		Name:    ContentNameTemplateCode,
		Example: "SMS_153055065",
		Escape:  "",
	},
	{
		Name:    ContentNameTemplateParam,
		Example: `{"code":"1111"}`,
		Escape:  "jsonescape",
	},
	{
		Name:    ContentNameSmsUpExtendCode,
		Example: `90999`,
		Escape:  "",
	},
	{
		Name:    ContentNameOutID,
		Example: `abcdefgh`,
		Escape:  "",
	},
}
