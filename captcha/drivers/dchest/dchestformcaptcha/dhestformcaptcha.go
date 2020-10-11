package dchestformcaptcha

import (
	"encoding/base64"
	"net/http"

	"github.com/herb-go/herbmodules/captcha"
	"github.com/herb-go/herbmodules/captcha/formcaptcha"
)

type Driver struct {
	captcha *captcha.Captcha
}

func NewDriver() *Driver {
	return &Driver{}
}

func (d *Driver) CreateContext(scene captcha.Scene, r *http.Request) (*formcaptcha.Context, error) {
	var err error
	ctx := formcaptcha.NewContext()
	ctx.Context = captcha.NewPlainContext()
	err = ctx.SetCaptchaData(captcha.ContextNameSubmited, []byte(r.Header.Get("captchasubmited")))
	if err != nil {
		return nil, err
	}
	err = ctx.SetCaptchaData(captcha.ContextNameTimestamp, []byte(r.Header.Get("captchatimestamp")))
	if err != nil {
		return nil, err
	}
	err = ctx.SetCaptchaData(captcha.ContextNameToken, []byte(r.Header.Get("captchatoken")))
	if err != nil {
		return nil, err
	}
	err = ctx.SetCaptchaData(captcha.ContextNameSign, []byte(r.Header.Get("captchasign")))
	if err != nil {
		return nil, err
	}
	return ctx, nil
}
func (d *Driver) GetFailMessage(scene captcha.Scene) (string, error) {
	return "captcha error", nil
}
func (d *Driver) RenderOutput(ctx *formcaptcha.Context) (interface{}, error) {
	image, err := ctx.GetCaptchaData(captcha.ContextNameImage)
	if err != nil {
		return nil, err
	}
	sign, err := ctx.GetCaptchaData(captcha.ContextNameSign)
	if err != nil {
		return nil, err
	}
	token, err := ctx.GetCaptchaData(captcha.ContextNameToken)
	if err != nil {
		return nil, err
	}
	timestamp, err := ctx.GetCaptchaData(captcha.ContextNameTimestamp)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"Image":     base64.StdEncoding.EncodeToString(image),
		"Sign":      string(sign),
		"Token":     string(token),
		"Timestamp": string(timestamp),
	}, nil
}

func (d *Driver) Captcha() *captcha.Captcha {
	return d.captcha
}
