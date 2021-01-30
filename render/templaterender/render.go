package templaterender

import (
	"time"

	"github.com/herb-go/herbtext"
	"github.com/herb-go/herbtext/texttemplate"
	"github.com/herb-go/herbtext/texttemplate/texttemplateset"
	"github.com/herb-go/notification"
)

type Renderer struct {
	Name                string
	Description         string
	Topic               string
	TTL                 time.Duration
	Delivery            string
	Constants           herbtext.Set
	Params              *texttemplate.Params
	TemplateSet         texttemplateset.TemplateSet
	SupportedDirectives []string
}

//Render render notification with given data
func (r *Renderer) Render(data map[string]string) (*notification.Notification, error) {
	c := notification.NewContent()
	herbtext.MergeSet(c, herbtext.Map(data))
	herbtext.MergeSet(c, r.Constants)
	ds, err := r.Params.Load(c)
	if err != nil {
		return nil, err
	}
	content, err := r.TemplateSet.Render(ds)
	if err != nil {
		return nil, err
	}
	n := notification.New()
	n.Header.Set(notification.HeaderNameTopic, r.Topic)
	ed := r.TTL
	if ed <= 0 {
		ed = notification.SuggestedNotificationTTL
	}
	n.ExpiredTime = time.Now().Add(ed).Unix()
	n.Delivery = r.Delivery
	herbtext.MergeSet(n.Content, content)
	return n, nil
}

//Supported return supported directives.
func (r *Renderer) Supported() (directives []string, err error) {
	return r.SupportedDirectives, nil
}
