package overseers

import (
	"github.com/herb-go/herb/middleware/httpinfo"
	"github.com/herb-go/herbmodules/httpinfomanager"
	"github.com/herb-go/worker"
)

var extractorworker httpinfo.Extractor
var extractorTeam = worker.GetWorkerTeam(&extractorworker)

func GetExtractorByID(id string) httpinfo.Extractor {
	w := worker.FindWorker(id)
	if w == nil {
		return nil
	}
	c, ok := w.Interface.(*httpinfo.Extractor)
	if ok == false || c == nil {
		return nil
	}
	return *c
}

var extractorfactoryworker httpinfomanager.ExtractorFactory
var extractorfactoryTeam = worker.GetWorkerTeam(&extractorfactoryworker)

func GetExtractorFactoryByID(id string) httpinfomanager.ExtractorFactory {
	w := worker.FindWorker(id)
	if w == nil {
		return nil
	}
	c, ok := w.Interface.(*httpinfomanager.ExtractorFactory)
	if ok == false || c == nil {
		return nil
	}
	return *c
}

var formaterworker httpinfo.Formatter
var FormatterTeam = worker.GetWorkerTeam(&formaterworker)

func GetFormatterByID(id string) httpinfo.Formatter {
	w := worker.FindWorker(id)
	if w == nil {
		return nil
	}
	c, ok := w.Interface.(*httpinfo.Formatter)
	if ok == false || c == nil {
		return nil
	}
	return *c
}

var formatterfactoryworker httpinfomanager.FormatterFactory
var FormatterFactoryTeam = worker.GetWorkerTeam(&formatterfactoryworker)

func GetFormatterFactoryByID(id string) httpinfomanager.FormatterFactory {
	w := worker.FindWorker(id)
	if w == nil {
		return nil
	}
	c, ok := w.Interface.(*httpinfomanager.FormatterFactory)
	if ok == false || c == nil {
		return nil
	}
	return *c
}
