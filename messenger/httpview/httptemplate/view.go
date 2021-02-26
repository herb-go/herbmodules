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
	Name        string
	Description string
	Config      *templateview.Config
}
type Template struct {
	TTLInSeconds    int64
	Constants       map[string]string
	Delivery        string
	Topic           string
	ContentTemplate map[string]string
	HeaderTemplate  map[string]string
	Required        []string
	Params          texttemplate.ParamDefinitions
	Engine          string
	Data            map[string]string
}

func (t *Template) config() *templateview.Config {
	return &templateview.Config{
		Params:          t.Params,
		Constants:       t.Constants,
		TTLInSeconds:    t.TTLInSeconds,
		ContentTemplate: t.ContentTemplate,
		HeaderTemplate:  t.HeaderTemplate,
		Topic:           t.Topic,
		Required:        t.Required,
		Engine:          t.Engine,
		Delivery:        t.Delivery,
	}
}
func (t *Template) Parse() (*templateview.View, error) {
	c := t.config()
	return c.Create()
}
func (t *Template) MustTOML() []byte {
	o := &Output{
		Name:        "VIEWNAME",
		Description: "VIEWDESCRPIPTION",
		Config:      t.config(),
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
	if t.Delivery == "" {
		return "Delivery", nil
	}
	if t.Engine == "" {
		return "Engine", nil
	}
	engine, err := texttemplate.GetEngine(t.Engine)
	if err != nil {
		return "Engine", nil
	}
	for k, v := range t.HeaderTemplate {
		_, err := engine.Parse(v, herbtext.DefaultEnvironment())
		if err != nil {
			return fmt.Sprintf("%s.%s", "HeaderTemplate", k), nil
		}
	}
	for k, v := range t.ContentTemplate {
		_, err := engine.Parse(v, herbtext.DefaultEnvironment())
		if err != nil {
			return fmt.Sprintf("%s.%s", "ContentTemplate", k), nil
		}
	}
	return "", nil
}
