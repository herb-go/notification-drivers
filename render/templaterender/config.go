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

type RendererConfig struct {
	Name         string
	Description  string
	Topic        string
	TTLInSeconds int64
	Delivery     string
	Engine       string
	Constants    map[string]string
	Params       texttemplate.ParamDefinitions
	Templates    map[string]string
}

func (c *RendererConfig) CreateRenderer() (notificationrender.Renderer, error) {
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
	r.TTL = time.Second * time.Duration(c.TTLInSeconds)
	r.Constants = notification.NewContent()
	herbtext.MergeSet(r.Constants, herbtext.Map(c.Constants))
	ts := notification.NewContent()
	herbtext.MergeSet(ts, herbtext.Map(c.Templates))
	env := herbtext.DefaultEnvironment()
	r.Params, err = c.Params.CreateParams(env)
	if err != nil {
		return nil, err
	}
	r.TemplateSet, err = texttemplateset.ParseWithEngineName(ts, c.Engine, env)
	if err != nil {
		return nil, err
	}
	return r, nil
}
