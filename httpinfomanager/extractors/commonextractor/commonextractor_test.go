package commonextractor_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/herb-go/herbconfig/loader"
	_ "github.com/herb-go/herbconfig/loader/drivers/jsonconfig"
	"github.com/herb-go/herbmodules/httpinfomanager/extractors/commonextractor"

	"github.com/herb-go/herb/middleware/router"
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

func TestExtractor(t *testing.T) {
	var data string
	var err error
	var req *http.Request
	httpinfomanager.Reset()
	defer httpinfomanager.Reset()
	commonextractor.RegsiterFactories()
	req = httptest.NewRequest("GET", "http://127.0.0.1", nil)
	data, err = extractString("test.notfound", nil, req)
	if data != "" || err == nil {
		t.Fatal(data, err)
	}
	commonconfig := &commonextractor.FieldConfig{
		Field: "test",
	}
	req = httptest.NewRequest("GET", "http://127.0.0.1", nil)
	req.Header.Set("test", "test")
	data, err = extractString("header", newLoader(commonconfig), req)
	if data != "test" || err != nil {
		t.Fatal(data, err)
	}
	req = httptest.NewRequest("GET", "http://127.0.0.1?test=test", nil)
	data, err = extractString("query", newLoader(commonconfig), req)
	if data != "test" || err != nil {
		t.Fatal(data, err)
	}
	req = httptest.NewRequest("POST", "http://127.0.0.1", bytes.NewBufferString("test=test"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	data, err = extractString("form", newLoader(commonconfig), req)
	if data != "test" || err != nil {
		t.Fatal(data, err)
	}
	req = httptest.NewRequest("GET", "http://127.0.0.1?test=test", nil)
	router.GetParams(req).Set("test", "test")
	data, err = extractString("router", newLoader(commonconfig), req)
	if data != "test" || err != nil {
		t.Fatal(data, err)
	}
	req = httptest.NewRequest("GET", "http://127.0.0.1?test=test", nil)
	data, err = extractString("fixed", newLoader(commonconfig), req)
	if data != "test" || err != nil {
		t.Fatal(data, err)
	}
	req = httptest.NewRequest("GET", "http://127.0.0.1?test=test", nil)
	reqcookie := &http.Cookie{
		Name:  "test",
		Value: "test",
	}
	req.AddCookie(reqcookie)
	data, err = extractString("cookie", newLoader(commonconfig), req)
	if data != "test" || err != nil {
		t.Fatal(data, err)
	}
	req = httptest.NewRequest("GET", "http://127.0.0.1", nil)
	data, err = extractString("cookie", newLoader(commonconfig), req)
	if data != "" || err != nil {
		t.Fatal(data, err)
	}
	req = httptest.NewRequest("GET", "http://127.0.0.1", nil)
	req.RemoteAddr = "127.0.0.1:8000"
	data, err = extractString("ip", newLoader(commonconfig), req)
	if data != "127.0.0.1" || err != nil {
		t.Fatal(data, err)
	}
	req = httptest.NewRequest("GET", "http://127.0.0.1", nil)
	data, err = extractString("method", newLoader(commonconfig), req)
	if data != "GET" || err != nil {
		t.Fatal(data, err)
	}
	req = httptest.NewRequest("GET", "http://127.0.0.1/test", nil)
	req.URL.RawPath = "/test"
	data, err = extractString("path", newLoader(commonconfig), req)
	if data != "/test" || err != nil {
		t.Fatal(data, err)
	}
	req = httptest.NewRequest("GET", "http://127.0.0.1/test", nil)
	data, err = extractString("host", newLoader(commonconfig), req)
	if data != "127.0.0.1" || err != nil {
		t.Fatal(data, err)
	}
	req = httptest.NewRequest("GET", "http://127.0.0.1/test", nil)
	req.SetBasicAuth("testuser", "testpassword")
	data, err = extractString("user", newLoader(commonconfig), req)
	if data != "testuser" || err != nil {
		t.Fatal(data, err)
	}
	req = httptest.NewRequest("GET", "http://127.0.0.1/test", nil)
	req.SetBasicAuth("testuser", "testpassword")
	data, err = extractString("password", newLoader(commonconfig), req)
	if data != "testpassword" || err != nil {
		t.Fatal(data, err)
	}

}
