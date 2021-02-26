package noitificationapidelivery

import (
	"encoding/json"

	"github.com/herb-go/fetcher"
	"github.com/herb-go/herbmodules/messenger/httpdelivery"
	"github.com/herb-go/notification"
	"github.com/herb-go/notification/notificationdelivery"
)

type Delivery struct {
	Preset          *fetcher.Preset
	Type            string
	Fields          []*notificationdelivery.Field
	RequiredContent []string
}

//DeliveryType Delivery type
func (d *Delivery) DeliveryType() string {
	return d.Type
}

//MustEscape delivery escape helper
func (d *Delivery) MustEscape(u string) string {
	return u
}

//ContentFields return content fields
//Return invalid fields and any error raised
func (d *Delivery) ContentFields() []*notificationdelivery.Field {
	return d.Fields
}

//CheckInvalidContent check if given content invalid
//Return invalid fields and any error raised
func (d *Delivery) CheckInvalidContent(c notification.Content) ([]string, error) {
	return notification.CheckRequiredContent(c, d.RequiredContent), nil
}

//Deliver send give content.
//Return delivery status and any receipt if returned,and any error if raised.
func (d *Delivery) Deliver(c notification.Content) (status notificationdelivery.DeliveryStatus, receipt string, err error) {
	result := &httpdelivery.DeliveryResult{}
	bs := []byte{}
	resp, err := d.Preset.FetchWithJSONBodyAndParse(c, fetcher.Should200(fetcher.AsBytes(&bs)))
	if err != nil {
		return notificationdelivery.DeliveryStatusFail, "", err
	}
	err = json.Unmarshal(bs, result)
	if err != nil {
		return notificationdelivery.DeliveryStatusFail, "", resp
	}
	return result.Status, result.Msg, nil
}

type Config struct {
	Server          fetcher.Server
	RequiredContent []string
	Fields          []*notificationdelivery.Field
	DeliveryType    string
}

var Factory = func(loader func(interface{}) error) (notificationdelivery.DeliveryDriver, error) {
	c := &Config{}
	err := loader(c)
	if err != nil {
		return nil, err
	}
	p, err := c.Server.CreatePreset()
	if err != nil {
		return nil, err
	}
	d := &Delivery{
		Preset:          p.With(fetcher.Method("POST")),
		RequiredContent: c.RequiredContent,
		Type:            c.DeliveryType,
		Fields:          c.Fields,
	}
	if d.Type == "" {
		d.Type = DeliveryType
	}
	return d, nil
}

func init() {
	notificationdelivery.Register(DeliveryType, Factory)
}
