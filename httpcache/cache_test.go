package httpcache

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/herb-go/datamodules/herbcache"
	"github.com/herb-go/datamodules/herbcache/cachepreset"
	"github.com/herb-go/datamodules/herbcache/kvengine"
	"github.com/herb-go/herb/identifier"
	_ "github.com/herb-go/herbdata-drivers/kvdb-drivers/freecachedb"
	"github.com/herb-go/herbdata/kvdb"
	_ "github.com/herb-go/herbdata/kvdb/commonkvdb"
)

var factory = func() *herbcache.Storage {
	s := herbcache.NewStorage()
	config := &kvengine.StorageConfig{
		Cache: &kvdb.Config{
			Driver: "freecache",
			Config: func(v interface{}) error {
				return json.Unmarshal([]byte(`{"Size":50000}`), v)
			},
		},
		VersionTTL: 3600,
		VersionStore: &kvdb.Config{
			Driver: "inmemory",
		},
	}

	err := config.ApplyTo(s)
	if err != nil {
		panic(err)
	}
	return s
}

func newTestPreset(ttl int64) *cachepreset.Preset {
	c := herbcache.New()
	s := factory()
	c = c.OverrideStorage(s)
	return cachepreset.New(cachepreset.Cache(c), cachepreset.TTL(ttl))

}

var content int

func testAction(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("ts", strconv.FormatInt(time.Now().UnixNano(), 10))
	w.WriteHeader(content)
	w.Write([]byte(r.Header.Get("test") + strconv.Itoa(content)))

}

var testid = identifier.IDFunc(func(r *http.Request) (string, error) {
	return r.Header.Get("test"), nil
})

func TestResponseCache(t *testing.T) {
	p := newTestPreset(3600)
	c := New().OverridePreset(p).OverrideIdentifier(testid)
	content = 200
	mux := http.NewServeMux()
	mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("rawts", strconv.FormatInt(time.Now().UnixNano(), 10))
		c.ServeMiddleware(w, r, testAction)
	})
	mux.HandleFunc("/raw", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("rawts", strconv.FormatInt(time.Now().UnixNano(), 10))
		testAction(w, r)
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	req, err := http.NewRequest("GET", server.URL+"/test", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("test", "test1")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	ts := resp.Header.Get("ts")
	rawts := resp.Header.Get("rawts")

	resp.Body.Close()

	if string(bs) != "test1"+"200" {
		t.Error(string(bs))
	}
	if resp.StatusCode != 200 {
		t.Error(resp.StatusCode)
	}

	req, err = http.NewRequest("GET", server.URL+"/raw", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("test", "test1")

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	bs, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()

	if string(bs) != "test1"+"200" {
		t.Error(string(bs))
	}
	if resp.StatusCode != 200 {
		t.Error(resp.StatusCode)
	}

	content = 404
	req, err = http.NewRequest("GET", server.URL+"/test", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("test", "test1")

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	bs, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()

	if string(bs) != "test1"+"200" {
		t.Fatal(string(bs))
	}
	if resp.Header.Get("ts") != ts {
		t.Error(resp.Header.Get("ts"))
	}
	if resp.Header.Get("rawts") == rawts {
		t.Error(resp.Header.Get("rawts"))
	}
	if resp.StatusCode != 200 {
		t.Error(resp.StatusCode)
	}

	req, err = http.NewRequest("GET", server.URL+"/raw", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("test", "test1")

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	bs, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()

	if string(bs) != "test1"+"404" {
		t.Error(string(bs))
	}
	if resp.Header.Get("ts") == ts {
		t.Error(resp.Header.Get("ts"))
	}
	if resp.Header.Get("rawts") == rawts {
		t.Error(resp.Header.Get("rawts"))
	}
	if resp.StatusCode != 404 {
		t.Error(resp.StatusCode)
	}

	req, err = http.NewRequest("GET", server.URL+"/test", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("test", "test2")

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	bs, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()

	if string(bs) != "test2"+"404" {
		t.Error(string(bs))
	}
	if resp.Header.Get("ts") == ts {
		t.Error(resp.Header.Get("ts"))
	}
	if resp.Header.Get("rawts") == rawts {
		t.Error(resp.Header.Get("rawts"))
	}
	if resp.StatusCode != 404 {
		t.Error(resp.StatusCode)
	}
	content = 500
	req, err = http.NewRequest("GET", server.URL+"/test", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("test", "test3")

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	bs, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()

	if string(bs) != "test3"+"500" {
		t.Error(string(bs))
	}
	if resp.StatusCode != 500 {
		t.Error(resp.StatusCode)
	}
	content = 403
	req, err = http.NewRequest("GET", server.URL+"/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("test", "test3")
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	bs, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()

	if string(bs) != "test3"+"403" {
		t.Error(string(bs))
	}
	if resp.StatusCode != 403 {
		t.Error(resp.StatusCode)
	}

}
