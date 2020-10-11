package httpinfomanager

import "github.com/herb-go/herb/middleware/httpinfo"

type FormatterFactory interface {
	CreateFormatter(loader func(interface{}) error) (httpinfo.Formatter, error)
}

type FormatterFactoryFunc func(loader func(interface{}) error) (httpinfo.Formatter, error)

func (f FormatterFactoryFunc) CreateFormatter(loader func(interface{}) error) (httpinfo.Formatter, error) {
	return f(loader)
}
