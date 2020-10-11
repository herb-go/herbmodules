package formcaptcha

import (
	"net/http"

	"github.com/herb-go/herbmodules/captcha"
)

type Context struct {
	captcha.Context
	Request *http.Request
	Output  interface{}
}

func NewContext() *Context {
	return &Context{}
}
