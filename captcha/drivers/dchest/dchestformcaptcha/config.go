package dchestformcaptcha

import (
	"github.com/herb-go/herbmodules/captcha"
	"github.com/herb-go/herbmodules/captcha/drivers/dchest"
	"github.com/herb-go/herbmodules/captcha/formcaptcha"
)

type Config struct {
	*dchest.Config
}

func Register() {
	formcaptcha.Register("dchest", func(loader func(v interface{}) error) (formcaptcha.Driver, error) {
		c := &Config{}
		err := loader(c)
		if err != nil {
			return nil, err
		}
		d := NewDriver()
		captchadriver, err := c.Config.Create()
		if err != nil {
			return nil, err
		}
		d.captcha = captcha.New()
		d.captcha.Driver = captchadriver
		return d, nil
	})
}

func init() {
	Register()
}
