package httptemplate

import (
	"bytes"
	"fmt"

	"github.com/BurntSushi/toml"

	"github.com/herb-go/herbtext"
	"github.com/herb-go/herbtext/texttemplate"
	"github.com/herb-go/notification-drivers/view/templateview"
)

type Output struct {
	Views []*View
}

type View struct {
	Name        string
	Description string
	Type        string
	Config      *templateview.Config
}
type TemplateConfig struct {
	TTLInSeconds    int64
	Delivery        string
	Topic           string
	ContentTemplate map[string]string
	HeaderTemplate  map[string]string
	Params          texttemplate.ParamDefinitions
	Engine          string
}
type Template struct {
	Config TemplateConfig
	Data   map[string]string
}

func (t *Template) toConfig() *templateview.Config {
	return &templateview.Config{
		Params:          t.Config.Params,
		TTLInSeconds:    t.Config.TTLInSeconds,
		ContentTemplate: t.Config.ContentTemplate,
		HeaderTemplate:  t.Config.HeaderTemplate,
		Topic:           t.Config.Topic,
		Engine:          t.Config.Engine,
		Delivery:        t.Config.Delivery,
	}
}
func (t *Template) Parse() (*templateview.View, error) {
	c := t.toConfig()
	return c.Create()
}
func (t *Template) MustTOML() []byte {
	o := &Output{
		Views: []*View{
			{
				Name:        "VIEWNAME",
				Description: "",
				Type:        "template",
				Config:      t.toConfig(),
			},
		},
	}
	buf := bytes.NewBuffer(nil)
	e := toml.NewEncoder(buf)
	err := e.Encode(o)
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
}
func (t *Template) Validate() (string, error) {
	if t.Config.Delivery == "" {
		return "Delivery", nil
	}
	if t.Config.Engine == "" {
		return "Engine", nil
	}
	engine, err := texttemplate.GetEngine(t.Config.Engine)
	if err != nil {
		return "Engine", nil
	}
	for k, v := range t.Config.HeaderTemplate {
		_, err := engine.Parse(v, herbtext.DefaultEnvironment())
		if err != nil {
			return fmt.Sprintf("%s.%s", "HeaderTemplate", k), nil
		}
	}
	for k, v := range t.Config.ContentTemplate {
		_, err := engine.Parse(v, herbtext.DefaultEnvironment())
		if err != nil {
			return fmt.Sprintf("%s.%s", "ContentTemplate", k), nil
		}
	}
	return "", nil
}
