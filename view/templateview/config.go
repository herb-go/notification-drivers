package templateview

import (
	"errors"
	"time"

	"github.com/herb-go/herbtext/texttemplate/texttemplateset"
	"github.com/herb-go/notification/notificationview"

	"github.com/herb-go/herbtext"

	"github.com/herb-go/notification"

	"github.com/herb-go/herbtext/texttemplate"
)

//Config view config struct
type Config struct {
	//Topic notifitaction topc
	Topic string
	//TTL notifitcation ttl in seconds.
	//Notification.SuggestedTTL will be used if ttl <=0.
	TTLInSeconds int64
	//Delivery notifiacation delivery
	Delivery string
	//Engine text template engine name
	Engine string
	//Params render params
	Params texttemplate.ParamDefinitions
	//HeaderTemplate content template set.
	HeaderTemplate map[string]string
	//ContentTemplate content template set.
	ContentTemplate map[string]string
}

//Create create renderer
func (c *Config) Create() (*View, error) {
	var err error
	if c.Delivery == "" {
		return nil, errors.New("templaterenderer: empty name")
	}
	v := &View{}
	v.Topic = c.Topic
	v.Delivery = c.Delivery
	if c.TTLInSeconds > 0 {
		v.TTL = time.Second * time.Duration(c.TTLInSeconds)
	} else {
		v.TTL = notification.SuggestedNotificationTTL
	}
	contenttemplate := notification.NewContent()
	herbtext.MergeSet(contenttemplate, herbtext.Map(c.ContentTemplate))
	headertemplate := notification.NewHeader()
	herbtext.MergeSet(headertemplate, herbtext.Map(c.HeaderTemplate))

	env := herbtext.DefaultEnvironment()
	v.Params, err = c.Params.CreateParams(env)
	if err != nil {
		return nil, err
	}
	eng, err := texttemplate.GetEngine(c.Engine)
	if err != nil {
		return nil, err
	}
	v.ContentTemplate, err = texttemplateset.ParseWith(contenttemplate, eng, env)
	if err != nil {
		return nil, err
	}
	v.HeaderTemplate, err = texttemplateset.ParseWith(headertemplate, eng, env)
	if err != nil {
		return nil, err
	}
	v.SupportedDirectives, err = eng.Supported(env)
	if err != nil {
		return nil, err
	}
	return v, nil
}

//Factory view factory
func Factory(loader func(v interface{}) error) (notificationview.View, error) {
	c := &Config{}
	err := loader(c)
	if err != nil {
		return nil, err
	}
	return c.Create()
}

func init() {
	notificationview.Register(DriverName, Factory)
}
