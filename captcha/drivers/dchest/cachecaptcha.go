package dchest

import (
	"strconv"
	"time"

	"github.com/herb-go/herbmodules/captcha/drivers/dchest/dchestimage"

	"github.com/herb-go/herbmodules/captcha"
	"github.com/herb-go/deprecated/cache"
)

type Driver struct {
	Cache          cache.Cacheable
	CachePrefix    string
	Secret         string
	TTL            time.Duration
	Width          int
	Height         int
	Wanted         *captcha.Wanted
	ImageRender    *dchestimage.Config
	Signer         func(string) (string, error)
	TokenGenerator func() (string, error)
}

func (d *Driver) Sign(token string, ts string, code string) (string, error) {
	return d.Signer(token + ts + code + d.Secret)
}
func (d *Driver) DoCaptcha(ctx captcha.Context) (bool, error) {
	token, err := ctx.GetCaptchaData(captcha.ContextNameToken)
	if err != nil {
		return false, err
	}
	if token == nil {
		return false, nil
	}
	submited, err := ctx.GetCaptchaData(captcha.ContextNameSubmited)
	if err != nil {
		return false, err
	}
	if submited == nil {
		return false, nil
	}
	sign, err := ctx.GetCaptchaData(captcha.ContextNameSign)
	if err != nil {
		return false, err
	}
	if sign == nil {
		return false, nil
	}
	ts, err := ctx.GetCaptchaData(captcha.ContextNameTimestamp)
	if err != nil {
		return false, err
	}
	if ts == nil {
		return false, nil
	}
	tsint, err := strconv.ParseInt(string(ts), 10, 64)
	if err != nil {
		return false, nil
	}
	if !time.Unix(tsint, 0).Add(d.TTL).After(time.Now()) {
		return false, nil
	}
	s, err := d.Sign(string(token), string(ts), string(submited))
	if err != nil {
		return false, err
	}
	result := (s == string(sign))
	if !result {
		return false, nil
	}
	_, err = d.Cache.GetBytesValue(d.CachePrefix + cache.KeyPrefix + string(token))
	if err != cache.ErrNotFound {
		return false, nil
	}
	err = d.Cache.SetBytesValue(d.CachePrefix+cache.KeyPrefix+string(token), []byte{1}, d.TTL)
	if err != nil {
		return false, err
	}
	return result, nil

}
func (d *Driver) Challenge(ctx captcha.Context) error {
	token, err := d.TokenGenerator()
	if err != nil {
		return err
	}
	ts := strconv.FormatInt(time.Now().Unix(), 10)
	wanted := d.Wanted.NewWantedBytes()
	sign, err := d.Sign(token, ts, string(wanted))
	if err != nil {
		return err
	}
	err = ctx.SetCaptchaData(captcha.ContextNameToken, []byte(token))
	if err != nil {
		return err
	}
	err = ctx.SetCaptchaData(captcha.ContextNameTimestamp, []byte(ts))
	if err != nil {
		return err
	}
	err = ctx.SetCaptchaData(captcha.ContextNameSign, []byte(sign))
	if err != nil {
		return err
	}
	indexs, err := d.Wanted.ToIndexBytes(wanted)
	if err != nil {
		return err
	}
	image, err := d.ImageRender.SavePNG("", indexs, d.Width, d.Height)
	if err != nil {
		return err
	}
	err = ctx.SetCaptchaData(captcha.ContextNameImage, image)
	if err != nil {
		return err
	}

	return nil
}
