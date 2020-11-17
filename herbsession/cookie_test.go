package herbsession

import (
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"

	"testing"

	"github.com/herb-go/herbdata"

	"github.com/herb-go/herb/middleware"
	"github.com/herb-go/herb/service/httpservice/httpcookie"
)

func newTestMux(s *Store) *http.ServeMux {
	var mux = http.NewServeMux()
	mux.Handle("/get", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session := s.RequestSession(r)
		value, err := session.Load("key")
		if err != nil {
			if err == herbdata.ErrNotFound {
				http.Error(w, http.StatusText(404), 404)
				return
			}
			http.Error(w, err.Error(), 500)
		}
		w.Write(value)
	}))
	mux.Handle("/set", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session := s.RequestSession(r)
		err := session.Store("key", []byte(r.URL.Query().Get("value")))
		if err != nil {
			http.Error(w, err.Error(), 500)
		}
		w.Write([]byte("ok"))
	}))
	mux.Handle("/delete", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session := s.RequestSession(r)
		err := session.Remove("key")
		if err != nil {
			http.Error(w, err.Error(), 500)
		}
		w.Write([]byte("ok"))
	}))
	mux.Handle("/revoke", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := s.RevokeRequestSession(r)
		if err != nil {
			http.Error(w, err.Error(), 500)
		}
		w.Write([]byte("ok"))
	}))
	mux.Handle("/regenerateid", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := s.RegenerateRequestSessionID(r)
		if err != nil {
			http.Error(w, err.Error(), 500)
		}
		w.Write([]byte("ok"))
	}))
	mux.Handle("/temporay", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session := s.RequestSession(r)
		session.MarkAsTemporay()
		w.Write([]byte("ok"))
	}))
	return mux
}

func TestCookie(t *testing.T) {
	var err error
	var s *Store
	var c *Cookie
	var value []byte
	var client *http.Client
	var jar *cookiejar.Jar
	var req *http.Request
	var resp *http.Response
	var u *url.URL
	var cookies []*http.Cookie
	var oldtoken string
	s = newTestStore()
	s.AutoStart = true
	c = &Cookie{
		httpcookie.Config{
			Name: "session",
		},
	}
	s.AddInstaller(c)
	var app = middleware.New(s.MustInstallByID(InstallerIDCookie))

	app.Handle(newTestMux(s))
	server := httptest.NewServer(app)
	defer server.Close()
	u, err = url.Parse(server.URL)
	if err != nil {
		panic(err)
	}
	jar, err = cookiejar.New(nil)
	if err != nil {
		panic(err)
	}
	client = &http.Client{}
	client.Jar = jar
	req, err = http.NewRequest("GET", server.URL+"/get", nil)
	if err != nil {
		panic(err)
	}
	resp, err = client.Do(req)
	if err != nil {
		panic(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 404 {
		t.Fatal()
	}
	cookies = jar.Cookies(u)
	if len(cookies) != 1 {
		t.Fatal(len(cookies))
	}
	req, err = http.NewRequest("GET", server.URL+"/set?value=test", nil)
	if err != nil {
		panic(err)
	}
	resp, err = client.Do(req)
	if err != nil {
		panic(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Fatal()
	}
	req, err = http.NewRequest("GET", server.URL+"/get", nil)
	if err != nil {
		panic(err)
	}
	resp, err = client.Do(req)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != 200 {
		t.Fatal()
	}
	value, err = ioutil.ReadAll(resp.Body)
	if string(value) != "test" {
		t.Fatal(string(value))
	}
	resp.Body.Close()

	req, err = http.NewRequest("GET", server.URL+"/delete", nil)
	if err != nil {
		panic(err)
	}
	resp, err = client.Do(req)
	if err != nil {
		panic(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Fatal()
	}
	req, err = http.NewRequest("GET", server.URL+"/get", nil)
	if err != nil {
		panic(err)
	}
	resp, err = client.Do(req)
	if err != nil {
		panic(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 404 {
		t.Fatal()
	}
	cookies = jar.Cookies(u)
	if len(cookies) != 1 {
		t.Fatal(len(cookies))
	}
	req, err = http.NewRequest("GET", server.URL+"/set?value=test", nil)
	if err != nil {
		panic(err)
	}
	resp, err = client.Do(req)
	if err != nil {
		panic(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Fatal()
	}
	req, err = http.NewRequest("GET", server.URL+"/get", nil)
	if err != nil {
		panic(err)
	}
	resp, err = client.Do(req)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != 200 {
		t.Fatal()
	}
	cookies = jar.Cookies(u)
	if len(cookies) != 1 {
		t.Fatal(len(cookies))
	}
	req, err = http.NewRequest("GET", server.URL+"/revoke", nil)
	if err != nil {
		panic(err)
	}
	resp, err = client.Do(req)
	if err != nil {
		panic(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Fatal()
	}
	cookies = jar.Cookies(u)
	if len(cookies) != 0 && cookies[0].Value != "" {
		t.Fatal(len(cookies), cookies[0].Value)
	}
	req, err = http.NewRequest("GET", server.URL+"/get", nil)
	if err != nil {
		panic(err)
	}
	resp, err = client.Do(req)
	if err != nil {
		panic(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 404 {
		t.Fatal()
	}
	req, err = http.NewRequest("GET", server.URL+"/set?value=test", nil)
	if err != nil {
		panic(err)
	}
	resp, err = client.Do(req)
	if err != nil {
		panic(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Fatal()
	}

	cookies = jar.Cookies(u)
	oldtoken = cookies[0].Value
	req, err = http.NewRequest("GET", server.URL+"/regenerateid", nil)
	if err != nil {
		panic(err)
	}
	resp, err = client.Do(req)
	if err != nil {
		panic(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Fatal()
	}
	cookies = jar.Cookies(u)
	if len(cookies) != 1 && cookies[0].Value != oldtoken {
		t.Fatal(len(cookies), cookies[0].Value)
	}
	req, err = http.NewRequest("GET", server.URL+"/get", nil)
	if err != nil {
		panic(err)
	}
	resp, err = client.Do(req)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != 200 {
		t.Fatal()
	}
	value, err = ioutil.ReadAll(resp.Body)
	if string(value) != "test" {
		t.Fatal(string(value))
	}
	resp.Body.Close()
}
