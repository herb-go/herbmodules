package formcaptcha

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/herb-go/herbmodules/captcha"
)

type ContextKey string

func (k ContextKey) SetErrMsg(r **http.Request, msg string) {
	ctx := (*r).Context()
	req := (*r).WithContext(context.WithValue(ctx, k, msg))
	*r = req
}

func (k ContextKey) GetErrMsg(r *http.Request) (string, bool) {
	val := r.Context().Value(k)
	if val == nil {
		return "", false
	}
	msg, ok := val.(string)
	return msg, ok
}

var DefaultContextKey = ContextKey("")

type Driver interface {
	CreateContext(scene captcha.Scene, r *http.Request) (*Context, error)
	GetFailMessage(scene captcha.Scene) (string, error)
	RenderOutput(*Context) (interface{}, error)
	Captcha() *captcha.Captcha
}

type FormCaptcha struct {
	Name       string
	ContextKey ContextKey
	Driver
}

func (c *FormCaptcha) GetErrMsg(r *http.Request) (string, bool) {
	return c.ContextKey.GetErrMsg(r)
}
func (c *FormCaptcha) NewCaptchaAction(scenefield string) http.HandlerFunc {
	if scenefield == "" {
		scenefield = "scene"
	}
	return func(w http.ResponseWriter, r *http.Request) {
		scene := r.URL.Query().Get(scenefield)
		ctx, err := c.Driver.CreateContext(captcha.Scene(scene), r)
		if err != nil {
			panic(err)
		}
		err = c.Captcha().Challenge(ctx)
		if err != nil {
			panic(err)
		}
		data, err := c.RenderOutput(ctx)
		if err != nil {
			panic(err)
		}
		output := NewOutput()
		output.Captcha = c.Name
		output.Trusted, err = ctx.CaptchaTrusted()
		if err != nil {
			panic(err)
		}
		output.Data = data
		bs, err := json.Marshal(output)
		if err != nil {
			panic(err)
		}
		_, err = w.Write(bs)
		if err != nil {
			panic(err)
		}
	}
}
func (c *FormCaptcha) Middleware(scene string) func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		ctx, err := c.Driver.CreateContext(captcha.Scene(scene), r)
		if err != nil {
			panic(err)
		}
		result, err := c.Captcha().DoCaptcha(ctx)
		if err != nil {
			panic(err)
		}
		if !result {
			msg, err := c.GetFailMessage(captcha.Scene(scene))
			if err != nil {
				panic(err)
			}
			c.ContextKey.SetErrMsg(&r, msg)
		}
		next(w, r)
	}
}

func New() *FormCaptcha {
	return &FormCaptcha{}
}
