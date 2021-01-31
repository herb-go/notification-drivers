package templaterender

import (
	"errors"
	"time"

	"github.com/herb-go/herbtext/texttemplate/texttemplateset"

	"github.com/herb-go/herbtext"

	"github.com/herb-go/notification"

	"github.com/herb-go/herbtext/texttemplate"
	"github.com/herb-go/notification/notificationrender"
)

//RendererConfig renderer config struct
type RendererConfig struct {
	//Name renderer name
	Name string
	//Description renderer description
	Description string
	//Topic notifitaction topc
	Topic string
	//TTL notifitcation ttl in seconds.
	//Notification.SuggestedTTL will be used if ttl <=0.
	TTLInSeconds int64
	//Delivery notifiacation delivery
	Delivery string
	//Engine text template engine name
	Engine string
	//Constants constatns will overwrite given values
	Constants map[string]string
	//Params render params
	Params texttemplate.ParamDefinitions
	//HeaderTemplate content template set.
	HeaderTemplate map[string]string
	//ContentTemplate content template set.
	ContentTemplate map[string]string
}

//Create create renderer
func (c *RendererConfig) Create() (*Renderer, error) {
	var err error
	if c.Name == "" {
		return nil, errors.New("templaterenderer: empty name")
	}
	if c.Delivery == "" {
		return nil, errors.New("templaterenderer: empty name")
	}
	r := &Renderer{}
	r.Name = c.Name
	r.Description = c.Description
	r.Topic = c.Topic
	r.Delivery = c.Delivery
	if c.TTLInSeconds > 0 {
		r.TTL = time.Second * time.Duration(c.TTLInSeconds)
	} else {
		r.TTL = notification.SuggestedNotificationTTL
	}
	r.Constants = notification.NewContent()
	herbtext.MergeSet(r.Constants, herbtext.Map(c.Constants))
	contenttemplate := notification.NewContent()
	herbtext.MergeSet(contenttemplate, herbtext.Map(c.ContentTemplate))
	headertemplate := notification.NewHeader()
	herbtext.MergeSet(headertemplate, herbtext.Map(c.HeaderTemplate))

	env := herbtext.DefaultEnvironment()
	r.Params, err = c.Params.CreateParams(env)
	if err != nil {
		return nil, err
	}
	eng, err := texttemplate.GetEngine(c.Engine)
	if err != nil {
		return nil, err
	}
	r.ContentTemplate, err = texttemplateset.ParseWith(contenttemplate, eng, env)
	if err != nil {
		return nil, err
	}
	r.HeaderTemplate, err = texttemplateset.ParseWith(headertemplate, eng, env)
	if err != nil {
		return nil, err
	}
	r.SupportedDirectives, err = eng.Supported(env)
	if err != nil {
		return nil, err
	}
	return r, nil
}

//CreateRenderCenter create render center with given renderer config list.
func CreateRenderCenter(configs []*RendererConfig) (notificationrender.RenderCenter, error) {
	c := notificationrender.NewRenderCenter()
	for _, v := range configs {
		r, err := v.Create()
		if err != nil {
			return nil, err
		}
		c.Set(r.Name, r)
	}
	return c, nil
}

//Templates templates config
type Templates struct {
	Templates []*RendererConfig
}

//CreateRenderCenter CreateRenderCenter create render center.
func (t *Templates) CreateRenderCenter() (notificationrender.RenderCenter, error) {
	return CreateRenderCenter(t.Templates)
}
