package workerextractor

import (
	"fmt"
	"net/http"

	"github.com/herb-go/herb-drivers/overseers/identifieroverseer"
	"github.com/herb-go/herb/middleware/httpinfo"
	"github.com/herb-go/deprecated/httpuser"
	"github.com/herb-go/herbmodules/httpinfomanager"
	"github.com/herb-go/herbmodules/httpinfomanager/overseers"
	"github.com/herb-go/worker"
)

type IdentifierExtractor struct {
	Identifier httpuser.Identifier
}

func (httpinfo *IdentifierExtractor) Extract(r *http.Request) ([]byte, error) {
	id, err := httpinfo.Identifier.IdentifyRequest(r)
	if err != nil {
		return nil, err
	}
	return []byte(id), nil
}

type WorkerConfig struct {
	ID string
}

var IdentifierExtractorFactory = httpinfomanager.ExtractorFactoryFunc(func(loader func(v interface{}) error) (httpinfo.Extractor, error) {
	c := &WorkerConfig{}
	err := loader(c)
	if err != nil {
		return nil, err
	}
	identifier := identifieroverseer.GetIdentifierByID(c.ID)
	if identifier == nil {
		return nil, fmt.Errorf("httpinfomanager: %w (%s)", worker.ErrWorkerNotFound, c.ID)
	}
	extractor := &IdentifierExtractor{Identifier: identifier}
	return extractor, nil
})

var HiredExtractorFactory = httpinfomanager.ExtractorFactoryFunc(func(loader func(v interface{}) error) (httpinfo.Extractor, error) {
	c := &WorkerConfig{}
	err := loader(c)
	if err != nil {
		return nil, err
	}
	reader := overseers.GetExtractorByID(c.ID)
	if reader == nil {
		return nil, fmt.Errorf("workerextractor: %w (%s)", worker.ErrWorkerNotFound, c.ID)
	}
	return reader, nil
})

var WorkerDefaultExtractorFactory = func(name string) (httpinfomanager.ExtractorFactory, error) {
	f := overseers.GetExtractorFactoryByID(name)
	if f == nil {
		return nil, fmt.Errorf("workerextractor: %w (%s)", worker.ErrWorkerNotFound, name)
	}
	return f, nil

}

func RegsiterFactories() {
	httpinfomanager.RegisterExtractorFactory("identifier", IdentifierExtractorFactory)
	httpinfomanager.RegisterExtractorFactory("hired", HiredExtractorFactory)
	httpinfomanager.SetDefaultExtractorFactory(WorkerDefaultExtractorFactory)

}
