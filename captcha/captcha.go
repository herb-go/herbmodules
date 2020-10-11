package captcha

type Driver interface {
	DoCaptcha(Context) (bool, error)
	Challenge(Context) error
}

type Captcha struct {
	Driver
}

func (c *Captcha) DoCaptcha(ctx Context) (bool, error) {
	result, err := ctx.CaptchaTrusted()
	if result == true || err != nil {
		return result, err
	}
	return c.Driver.DoCaptcha(ctx)
}

func New() *Captcha {
	return &Captcha{}
}
