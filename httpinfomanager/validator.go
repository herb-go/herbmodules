package httpinfomanager

import (
	"github.com/herb-go/herb/middleware/httpinfo"
)

type ValidatorFactory interface {
	CreateValidator(loader func(interface{}) error) (httpinfo.Validator, error)
}

type ValidatorFactoryFunc func(loader func(interface{}) error) (httpinfo.Validator, error)

func (f ValidatorFactoryFunc) CreateValidator(loader func(interface{}) error) (httpinfo.Validator, error) {
	return f(loader)
}
