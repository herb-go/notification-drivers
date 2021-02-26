package tencentcloudsmsdelivery

import "github.com/herb-go/notification/notificationdelivery"

var DeliveryType = "tencentcloudsms"

var ContentNameTemplateID = "templateid"
var ContentNameSign = "sign"
var ContentNamePhoneNumber = "phonenumber"
var ContentNameTemplateParam = "templateparam"
var ContentNameSessionContext = "sessioncontext"
var ContentNameExtendCode = "extendcode"
var ContentNameSenderID = "senderid"

var RequeiredContent = []string{ContentNameTemplateID, ContentNamePhoneNumber}

var Fields = []*notificationdelivery.Field{
	{
		Name:    ContentNameTemplateID,
		Example: "TEMPLATEID",
		Escape:  "",
	},
	{
		Name:    ContentNameSign,
		Example: "COMAPNYSIGN",
		Escape:  "",
	},
	{
		Name:    ContentNamePhoneNumber,
		Example: "+8613500000000,+8613500000000",
		Escape:  "commaescape",
	},
	{
		Name:    ContentNameTemplateParam,
		Example: `12345,23456`,
		Escape:  "commaescape",
	},
	{
		Name:    ContentNameSessionContext,
		Example: ``,
		Escape:  "",
	},
	{
		Name:    ContentNameExtendCode,
		Example: ``,
		Escape:  "",
	},
	{
		Name:    ContentNameSenderID,
		Example: ``,
		Escape:  "",
	},
}
