package loggerdelivery

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/herb-go/logger"
	"github.com/herb-go/notification"
	"github.com/herb-go/notification/notificationdelivery"
)

type Delivery struct {
	logger       *logger.Logger
	delayMin     time.Duration
	delayMax     time.Duration
	failPercent  int
	errorPercent int
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
	delay := rand.Int63n(int64(d.delayMax - d.delayMin))
	time.Sleep(time.Duration(delay))
	if rand.Int31n(100) <= int32(d.failPercent) {
		return notificationdelivery.DeliveryStatusFail, "fail", nil
	}
	if rand.Int31n(100) <= int32(d.errorPercent) {
		return notificationdelivery.DeliveryStatusFail, "", errors.New("random error")
	}

	d.logger.Log(fmt.Sprintf("mock: %s", string(bs)))
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
	Logger           string
	DelayMinDuration string
	DelayMaxDuration string
	FailPercent      int
	ErrorPercent     int
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
		logger:       l,
		failPercent:  c.FailPercent,
		errorPercent: c.ErrorPercent,
	}
	if c.DelayMinDuration != "" {
		dmin, err := time.ParseDuration(c.DelayMinDuration)
		if err != nil {
			return nil, err
		}
		d.delayMin = dmin
	}
	if c.DelayMaxDuration != "" {
		dmax, err := time.ParseDuration(c.DelayMaxDuration)
		if err != nil {
			return nil, err
		}
		d.delayMax = dmax
	}
	if d.delayMax < d.delayMin {
		d.delayMax = d.delayMin
	}
	return d, nil
}

func init() {
	notificationdelivery.Register(DeliveryType, Factory)
}
