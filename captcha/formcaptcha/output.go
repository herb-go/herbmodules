package formcaptcha

type Mode string

type Output struct {
	Captcha string
	Scene   string
	Trusted bool
	Data    interface{}
}

func NewOutput() *Output {
	return &Output{}
}
