package templaterender

import (
	"time"

	"github.com/herb-go/herbtext"
	"github.com/herb-go/herbtext/texttemplate"
	"github.com/herb-go/herbtext/texttemplate/texttemplateset"
	"github.com/herb-go/notification"
)

//Renderer renderer struct
type Renderer struct {
	//Name renderer name
	Name string
	//Description renderer description
	Description string
	//Topic notifitaction topc
	Topic string
	//TTL notifitcation ttl
	TTL time.Duration
	//Delivery notifiacation delivery
	Delivery string
	//Required fields
	Required []string
	//Constants constatns will overwrite given values
	Constants herbtext.Set
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
func (r *Renderer) Render(message notification.Message) (*notification.Notification, error) {
	for k := range r.Required {
		if message[r.Required[k]] == "" {
			return nil, notification.NewRequiredContentError([]string{r.Required[k]})
		}
	}
	c := notification.NewContent()
	herbtext.MergeSet(c, message)
	herbtext.MergeSet(c, r.Constants)
	ds, err := r.Params.Load(c)
	if err != nil {
		return nil, err
	}
	content, err := r.ContentTemplate.Render(ds)
	if err != nil {
		return nil, err
	}
	header, err := r.HeaderTemplate.Render(ds)
	if err != nil {
		return nil, err
	}
	n := notification.New()
	herbtext.MergeSet(n.Header, header)
	n.Header.Set(notification.HeaderNameTopic, r.Topic)
	now := time.Now()
	n.CreatedTime = now.Unix()
	n.ExpiredTime = now.Add(r.TTL).Unix()
	n.Delivery = r.Delivery
	herbtext.MergeSet(n.Content, content)
	return n, nil
}

//Supported return supported directives.
func (r *Renderer) Supported() (directives []string, err error) {
	return r.SupportedDirectives, nil
}
