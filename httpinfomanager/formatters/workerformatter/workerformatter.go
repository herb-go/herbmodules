package workerformatter

import (
	"fmt"

	"github.com/herb-go/herb/middleware/httpinfo"
	"github.com/herb-go/herbmodules/httpinfomanager"
	"github.com/herb-go/herbmodules/httpinfomanager/overseers"
	"github.com/herb-go/worker"
)

type WorkerConfig struct {
	ID string
}

var HiredFormatterFactory = httpinfomanager.FormatterFactoryFunc(func(loader func(v interface{}) error) (httpinfo.Formatter, error) {
	c := &WorkerConfig{}
	err := loader(c)
	if err != nil {
		return nil, err
	}

	formater := overseers.GetFormatterByID(c.ID)
	if formater == nil {
		return nil, fmt.Errorf("workerformatter: %w (%s)", worker.ErrWorkerNotFound, c.ID)
	}
	return formater, nil
})

var WorkerDefaultFormatterFactory = func(name string) (httpinfomanager.FormatterFactory, error) {
	f := overseers.GetFormatterFactoryByID(name)
	if f == nil {
		return nil, fmt.Errorf("workerformatter: %w (%s)", worker.ErrWorkerNotFound, name)
	}
	return f, nil

}

func RegisterFactories() {
	httpinfomanager.RegisterFormatterFactory("hired", HiredFormatterFactory)
	httpinfomanager.SetDefaultFormatterFactory(WorkerDefaultFormatterFactory)

}
