package loggerdelivery

import (
	"encoding/json"
	"fmt"

	"github.com/herb-go/logger"
	"github.com/herb-go/notification"
	"github.com/herb-go/notification/notificationdelivery"
)

type Delivery struct {
	logger *logger.Logger
}

func (d *Delivery) DeliveryType() string {
	return DeliveryType
}

//CheckInvalidContent check if given content invalid
//Return invalid fields and any error raised
func (d *Delivery) CheckInvalidContent(notification.Content) ([]string, error) {
	return []string{}, nil
}
func (d *Delivery) Deliver(c notification.Content) (notificationdelivery.DeliveryStatus, string, error) {

	bs, err := json.Marshal(c)
	if err != nil {
		return notificationdelivery.DeliveryStatusAbort, "", err
	}
	d.logger.Log(fmt.Sprintf("loggerdelivery: %s", string(bs)))
	return notificationdelivery.DeliveryStatusSuccess, "", nil

}

func (d *Delivery) MustEscape(unescaped string) string {
	return unescaped
}

//ContentFields return content fields
//Return invalid fields and any error raised
func (d Delivery) ContentFields() []*notificationdelivery.Field {
	return nil
}

type Config struct {
	Logger string
}

var Factory = func(loader func(interface{}) error) (notificationdelivery.DeliveryDriver, error) {
	c := &Config{}
	err := loader(c)
	if err != nil {
		return nil, err
	}
	l := logger.GetBuiltinLogger(c.Logger)
	if l == nil {
		return nil, fmt.Errorf("builtin logger [%s] not found", c.Logger)
	}
	d := &Delivery{
		logger: l,
	}
	return d, nil
}

func init() {
	notificationdelivery.Register(DeliveryType, Factory)
}
