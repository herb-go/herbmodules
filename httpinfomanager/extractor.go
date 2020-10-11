package httpinfomanager

import "github.com/herb-go/herb/middleware/httpinfo"

type ExtractorFactory interface {
	CreateExtractor(loader func(interface{}) error) (httpinfo.Extractor, error)
}

type ExtractorFactoryFunc func(loader func(interface{}) error) (httpinfo.Extractor, error)

func (f ExtractorFactoryFunc) CreateExtractor(loader func(interface{}) error) (httpinfo.Extractor, error) {
	return f(loader)
}
