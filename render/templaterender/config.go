package templaterender

import "github.com/herb-go/herbtext/texttemplate"

type RenderConfig struct {
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
