package workerextractor_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/herb-go/herb/middleware/httpinfo"
	"github.com/herb-go/deprecated/httpuser"

	"github.com/herb-go/herbconfig/loader"
	_ "github.com/herb-go/herbconfig/loader/drivers/jsonconfig"
	"github.com/herb-go/herbmodules/httpinfomanager/extractors/workerextractor"
	"github.com/herb-go/worker"

	"github.com/herb-go/herbmodules/httpinfomanager"
)

func newLoader(v interface{}) func(interface{}) error {
	bytes, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return loader.NewLoader("json", bytes)
}

func extract(name string, loader func(interface{}) error, r *http.Request) ([]byte, error) {
	f, err := httpinfomanager.GetExtractorFactory(name)
	if err != nil {
		return nil, err
	}
	extractor, err := f.CreateExtractor(loader)
	if err != nil {
		return nil, err
	}
	return extractor.Extract(r)
}

func extractString(name string, loader func(interface{}) error, r *http.Request) (string, error) {
	b, err := extract(name, loader, r)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

type testidentifier string

func (i testidentifier) IdentifyRequest(r *http.Request) (string, error) {
	return string(i), nil
}

func TestExtractor(t *testing.T) {
	var data string
	var err error
	var req *http.Request
	httpinfomanager.Reset()
	defer httpinfomanager.Reset()
	workerextractor.RegsiterFactories()
	identifier := httpuser.Identifier(testidentifier("test"))
	worker.Hire("test.identifier", &identifier)
	identifierconfig := &workerextractor.WorkerConfig{ID: "test.identifier"}
	req = httptest.NewRequest("GET", "http://127.0.0.1/test", nil)
	data, err = extractString("test.notfound", nil, req)
	if data != "" || err == nil {
		t.Fatal(data, err)
	}
	data, err = extractString("identifier", newLoader(identifierconfig), req)
	if data != "test" || err != nil {
		t.Fatal(data, err)
	}
	workernotfoundconfig := &workerextractor.WorkerConfig{ID: "test.notfound"}
	req = httptest.NewRequest("GET", "http://127.0.0.1/test", nil)
	data, err = extractString("identifier", newLoader(workernotfoundconfig), req)
	if data != "" || err == nil {
		t.Fatal(data, err)
	}
	e := httpinfo.Extractor(httpinfo.ExtractorFunc(func(r *http.Request) ([]byte, error) {
		return []byte("test"), nil
	}))
	worker.Hire("test.extractor", &e)
	hiredconfig := &workerextractor.WorkerConfig{ID: "test.extractor"}
	data, err = extractString("hired", newLoader(hiredconfig), req)
	if data != "test" || err != nil {
		t.Fatal(data, err)
	}
	data, err = extractString("hired", newLoader(workernotfoundconfig), req)
	if data != "" || err == nil {
		t.Fatal(data, err)
	}
	factory := httpinfomanager.ExtractorFactory(httpinfomanager.ExtractorFactoryFunc(func(loader func(interface{}) error) (httpinfo.Extractor, error) {
		return e, nil
	}))
	worker.Hire("test.testfactory", &factory)
	req = httptest.NewRequest("GET", "http://127.0.0.1/test", nil)
	data, err = extractString("test.testfactory", newLoader(nil), req)
	if data != "test" || err != nil {
		t.Fatal(data, err)
	}
}
