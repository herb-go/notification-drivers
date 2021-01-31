package templaterender_test

import (
	"testing"
	"time"

	_ "github.com/herb-go/herbtext-drivers/engine/handlebars"
	"github.com/herb-go/notification"

	_ "github.com/herb-go/herbtext-drivers/commonenvironment"

	"github.com/herb-go/herbtext/texttemplate"
	"github.com/herb-go/notification-drivers/render/templaterender"
)

func TestTemplate(t *testing.T) {
	config := &templaterender.RendererConfig{
		Name:         "test",
		Description:  "test description",
		Topic:        "testtopic",
		TTLInSeconds: 3600,
		Delivery:     "testdelivery",
		Engine:       "handlebars",
		Constants: map[string]string{
			"constant":  "constant",
			"constant2": "constant2",
		},
		Params: texttemplate.ParamDefinitions{
			{
				ParamConfig: texttemplate.ParamConfig{
					Source: "testheader",
				},
			},
			{
				ParamConfig: texttemplate.ParamConfig{
					Source: "testdelivery",
				},
			},
			{
				ParamConfig: texttemplate.ParamConfig{
					Source: "constant",
				},
			},
			{
				ParamConfig: texttemplate.ParamConfig{
					Source: "constant2",
				},
			},
			{
				ParamConfig: texttemplate.ParamConfig{
					Source: "testkey",
				},
			},
			{
				ParamConfig: texttemplate.ParamConfig{
					Source: "testkey2",
				},
			},
		},
		HeaderTemplate: map[string]string{
			"testheader": "{{{testheader}}}",
			"topic":      "{{{testtopic}}}",
		},
		ContentTemplate: map[string]string{
			"testconstant":  "{{{constant}}}",
			"testconstant2": "{{{constant2}}}",
			"testkey":       "{{{testkey}}}",
			"testkey2":      "{{{testkey2}}}",
		},
	}
	rc, err := templaterender.CreateRenderCenter([]*templaterender.RendererConfig{config})
	if err != nil {
		t.Fatal(rc, err)
	}
	r, err := rc.Get("test")
	if err != nil {
		t.Fatal(r)
	}
	supported, err := r.Supported()
	if err != nil || len(supported) == 0 {
		t.Fatal(supported)
	}
	data := map[string]string{
		"testheader": "testheadervalue",
		"testtopic":  "testtopicvalue",
		"constant":   "constantvalue",
		"testkey":    "testkeyvalue",
		"testkey2":   "testkey2value",
	}
	n, err := r.Render(data)
	if err != nil || n == nil {
		t.Fatal(n, err)
	}
	if n.Delivery != "testdelivery" {
		t.Fatal(n)
	}
	if n.Header.Get(notification.HeaderNameTopic) != "testtopic" {
		t.Fatal(n)
	}
	if n.Content.Get("testconstant") != "constant" || n.Content.Get("testconstant2") != "constant2" {
		t.Fatal(n.Content)
	}
	if n.Header.Get("testheader") != "testheadervalue" {
		t.Fatal(n)
	}
	if n.CreatedTime > time.Now().Add(time.Minute).Unix() || n.CreatedTime < time.Now().Add(-time.Minute).Unix() {
		t.Fatal(n)
	}
	if n.ExpiredTime-n.CreatedTime != 3600 {
		t.Fatal(n.ExpiredTime - n.CreatedTime)
	}
}
