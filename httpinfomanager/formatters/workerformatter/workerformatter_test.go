package workerformatter_test

import (
	"encoding/json"
	"testing"

	"github.com/herb-go/herb/middleware/httpinfo"
	"github.com/herb-go/herbconfig/loader"
	_ "github.com/herb-go/herbconfig/loader/drivers/jsonconfig"
	"github.com/herb-go/herbmodules/httpinfomanager"
	"github.com/herb-go/herbmodules/httpinfomanager/formatters/workerformatter"
	"github.com/herb-go/worker"
)

func newLoader(v interface{}) func(interface{}) error {
	bytes, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return loader.NewLoader("json", bytes)
}
func format(name string, loader func(interface{}) error, data []byte) ([]byte, bool, error) {
	f, err := httpinfomanager.GetFormatterFactory(name)
	if err != nil {
		return nil, false, err
	}
	formatter, err := f.CreateFormatter(loader)
	if err != nil {
		return nil, false, err
	}
	return formatter.Format(data)
}
func formatString(name string, loader func(interface{}) error, datastr string) (string, bool, error) {
	data, ok, err := format(name, loader, []byte(datastr))
	return string(data), ok, err
}
func TestFormatter(t *testing.T) {
	// var data []byte
	var datastr string
	var ok bool
	var err error

	worker.Reset()
	defer worker.Reset()
	httpinfomanager.Reset()
	defer httpinfomanager.Reset()
	workerformatter.RegisterFactories()
	datastr, ok, err = formatString("test.notfound", nil, "abc")
	if datastr != "" || ok != false || err == nil {
		t.Fatal(datastr, ok, err)
	}
	formatter := httpinfo.Formatter(httpinfo.FormatterFunc(func(data []byte) ([]byte, bool, error) {
		return data, true, nil
	}))

	worker.Hire("test.formatter", &formatter)
	workerconfig := &workerformatter.WorkerConfig{
		ID: "test.formatter",
	}
	datastr, ok, err = formatString("hired", newLoader(workerconfig), "12345")
	if datastr != "12345" || ok != true || err != nil {
		t.Fatal(datastr, ok, err)
	}
	workernotfoundconfig := &workerformatter.WorkerConfig{
		ID: "test.notfound",
	}
	datastr, ok, err = formatString("hired", newLoader(workernotfoundconfig), "12345")
	if datastr != "" || ok != false || err == nil {
		t.Fatal(datastr, ok, err)
	}
	factory := httpinfomanager.FormatterFactory(httpinfomanager.FormatterFactoryFunc(func(loader func(interface{}) error) (httpinfo.Formatter, error) {
		return formatter, nil
	}))
	worker.Hire("test.formatterfactory", &factory)
	datastr, ok, err = formatString("test.formatterfactory", nil, "123-45")
	if datastr != "123-45" || ok != true || err != nil {
		t.Fatal(datastr, ok, err)
	}
}
