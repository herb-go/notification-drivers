package templateview

import (
	"time"

	"github.com/herb-go/herbtext"
	"github.com/herb-go/herbtext/texttemplate"
	"github.com/herb-go/herbtext/texttemplate/texttemplateset"
	"github.com/herb-go/notification"
	"github.com/herb-go/notification/notificationview"
)

//DriverName driver name
const DriverName = "template"

//View view struct
type View struct {
	//Topic notifitaction topc
	Topic string
	//TTL notifitcation ttl
	TTL time.Duration
	//Delivery notifiacation delivery
	Delivery string
	//Params render params
	Params *texttemplate.Params
	//HeaderTemplate header template set.
	HeaderTemplate texttemplateset.TemplateSet
	//ContentTemplate content template set.
	ContentTemplate texttemplateset.TemplateSet
	//SupportedDirectives renderer supported directives.
	SupportedDirectives []string
}

//Render render notification with given data
func (v *View) Render(message notificationview.Message) (*notification.Notification, error) {
	m := notificationview.NewMessage()
	herbtext.MergeSet(m, message)
	ds, err := v.Params.Load(m)
	if err != nil {
		if param := texttemplate.GetParamMissedErrorName(err); param != "" {
			return nil, notification.NewRequiredContentError([]string{param})
		}
		return nil, err
	}
	content, err := v.ContentTemplate.Render(ds)
	if err != nil {
		return nil, err
	}
	header, err := v.HeaderTemplate.Render(ds)
	if err != nil {
		return nil, err
	}
	n := notification.New()
	herbtext.MergeSet(n.Header, header)
	n.Header.Set(notification.HeaderNameTopic, v.Topic)
	now := time.Now()
	n.CreatedTime = now.Unix()
	n.ExpiredTime = now.Add(v.TTL).Unix()
	n.Delivery = v.Delivery
	herbtext.MergeSet(n.Content, content)
	return n, nil
}
