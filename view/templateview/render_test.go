package templateview_test

import (
	"testing"
	"time"

	"github.com/herb-go/notification/notificationview"

	_ "github.com/herb-go/herbtext-drivers/engine/handlebars"
	"github.com/herb-go/notification"
	"github.com/herb-go/notification-drivers/view/templateview"

	_ "github.com/herb-go/herbtext-drivers/commonenvironment"
	"github.com/herb-go/herbtext/texttemplate"
)

func newTestConfig() *templateview.Config {
	return &templateview.Config{
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
}
func TestTemplate(t *testing.T) {
	config := newTestConfig()
	viewconfig := &notificationview.ViewConfig{
		Type: templateview.DriverName,
		Config: func(v interface{}) error {
			c := v.(*templateview.Config)
			*c = *config
			return nil
		},
	}
	v, err := viewconfig.CreateView()
	if err != nil {
		t.Fatal(v, err)
	}
	if len(v.(*templateview.View).SupportedDirectives) == 0 {
		t.Fatal(v)
	}
	data := map[string]string{
		"testheader": "testheadervalue",
		"testtopic":  "testtopicvalue",
		"constant":   "constantvalue",
		"testkey":    "testkeyvalue",
		"testkey2":   "testkey2value",
	}
	n, err := v.Render(data)
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

func TestConfig(t *testing.T) {
	c := newTestConfig()
	c.Delivery = ""
	v, err := c.Create()
	if err == nil || v != nil {
		t.Fatal(err)
	}
	c = newTestConfig()
	c.TTLInSeconds = 0
	v, err = c.Create()
	if err != nil || v.TTL != notification.SuggestedNotificationTTL {
		t.Fatal(err)
	}
	c = newTestConfig()
	c.TTLInSeconds = -1
	v, err = c.Create()
	if err != nil || v.TTL != notification.SuggestedNotificationTTL {
		t.Fatal(err)
	}
}

func TestRequired(t *testing.T) {
	c := newTestConfig()
	c.Required = []string{"required"}
	v, err := c.Create()
	if err != nil || v == nil {
		t.Fatal(err)
	}
	data := map[string]string{
		"testheader": "testheadervalue",
		"testtopic":  "testtopicvalue",
		"constant":   "constantvalue",
		"testkey":    "testkeyvalue",
		"testkey2":   "testkey2value",
	}
	n, err := v.Render(data)
	if err == nil || !notification.IsInvalidContentError(err) {
		t.Fatal(n, err)
	}
	data["required"] = "required"
	n, err = v.Render(data)
	if err != nil || n == nil {
		t.Fatal(n, err)
	}
}
