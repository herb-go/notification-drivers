package kvreceiptstore

import (
	"strconv"
	"strings"

	"github.com/herb-go/notification"
	"github.com/herb-go/notification/notificationdelivery"
	"github.com/herb-go/notification/notificationdelivery/notificationqueue"
)

var SupportedConditions []string

func NewFilter() *Filter {
	return &Filter{
		NotificationFilter: notification.NewFilter(),
	}
}

type Filter struct {
	Status             *notificationdelivery.DeliveryStatus
	InMessage          string
	NotificationFilter *notification.PlainFilter
}

//FilterReceipt filter receipt with given context
//Return if Receipt is valid
func (f *Filter) FilterReceipt(r *notificationqueue.Receipt, ctx *notification.ConditionContext) (bool, error) {
	if f.Status != nil {
		if r.Status != *f.Status {
			return false, nil
		}
	}
	if f.InMessage != "" {
		if !strings.Contains(r.Message, f.InMessage) {
			return false, nil
		}
	}
	return f.NotificationFilter.FilterNotification(r.Notification, ctx)
}

//ApplyCondition apply search condition to filter
//ErrConditionNotSupported should be returned if condition keyword is not supported
func (f *Filter) ApplyCondition(cond *notification.Condition) error {
	switch cond.Keyword {
	case "status":
		i, err := strconv.ParseInt(cond.Value, 10, 64)
		if err != nil {
			return notification.NewErrInvalidConditionValue(cond.Value)
		}
		st := notificationdelivery.DeliveryStatus(i)
		f.Status = &st
		return nil
	case "inmessage":
		f.InMessage = cond.Value
		return nil
	}
	return f.NotificationFilter.ApplyCondition((cond))
}

//ApplyToFilter apply condiitons to filter.
func ApplyToFilter(f *Filter, conds []*notification.Condition) error {
	for k := range conds {
		err := f.ApplyCondition(conds[k])
		if err != nil {
			return err
		}
	}
	return nil
}

func init() {
	SupportedConditions = make([]string, len(notification.PlainFilterSupportedConditions))
	copy(SupportedConditions, notification.PlainFilterSupportedConditions)
	SupportedConditions = append(SupportedConditions, "stauts", "inmessage")
}
