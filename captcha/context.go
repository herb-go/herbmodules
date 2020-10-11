package captcha

type Context interface {
	SetCaptchaData(ContextName, ContextData) error
	GetCaptchaData(ContextName) (ContextData, error)
	CaptchaScene() (Scene, error)
	CaptchaTrusted() (bool, error)
}

type PlainContext struct {
	Scene   Scene
	Trusted bool
	*Collection
}

func (c *PlainContext) CaptchaScene() (Scene, error) {
	return c.Scene, nil
}
func (c *PlainContext) CaptchaTrusted() (bool, error) {
	return c.Trusted, nil
}

func NewPlainContext() *PlainContext {
	return &PlainContext{
		Scene:      DefaultScene,
		Collection: NewCollection(),
	}
}
