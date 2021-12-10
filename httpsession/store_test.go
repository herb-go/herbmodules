package httpsession

import (
	"bytes"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/herb-go/herb/middleware"
	"github.com/herb-go/herb/service/httpservice/httpcookie"
	"github.com/herb-go/herbdata"
)

func newTestStore() *Store {
	return &Store{
		Name:        "test",
		MaxLifetime: 1200,
		AutoStart:   false,
		Engine:      newTestAESEngine(),
	}
}
func TestStore(t *testing.T) {
	var s *Store
	var err error
	var session *Session
	var newsession *Session
	var data []byte
	var token string
	var e = newTestAESEngine()
	s = newTestStore()
	s.Timeout = 100
	s.MaxLifetime = 2
	s.Engine = e
	token, session, err = s.StartSession()
	if err != nil || session == nil {
		t.Fatal(err, token, session)
	}
	if session.Token() != token {
		t.Fatal()
	}
	err = session.Store("test", []byte("testdata"))
	if err != nil {
		t.Fatal()
	}
	err = s.SaveSession(session)
	if err != nil {
		t.Fatal()
	}
	time.Sleep(1 * time.Second)
	newsession, err = s.LoadSession(session.Token())
	if err != nil {
		t.Fatal()
	}
	if newsession.LoadedFrom() != session.Token() {
		t.Fatal(newsession.LoadedFrom())
	}
	data, err = newsession.Load("test")
	if err != nil {
		t.Fatal()
	}
	if bytes.Compare(data, []byte("testdata")) != 0 {
		t.Fatal()
	}
	err = s.RevokeSession(newsession.Token())
	if err != nil {
		t.Fatal()
	}
	token, session, err = s.StartSession()
	if err != nil || session == nil {
		t.Fatal(err, token, session)
	}
	err = session.Store("test", []byte("testdata"))
	if err != nil {
		t.Fatal()
	}
	err = s.SaveSession(session)
	if err != nil {
		t.Fatal()
	}
	_, err = s.LoadSession(session.Token())
	if err != nil {
		t.Fatal()
	}
	time.Sleep(3 * time.Second)
	_, err = s.LoadSession(session.Token())
	if err != herbdata.ErrNotFound {
		t.Fatal()
	}
}

func TestTemporay(t *testing.T) {
	var err error
	var s *Store
	var c *Cookie
	var client *http.Client
	var jar *cookiejar.Jar
	var req *http.Request
	var resp *http.Response
	var cookies []*http.Cookie
	s = newTestStore()
	c = &Cookie{
		httpcookie.Config{
			Name: "session",
		},
	}
	s.Installer = c
	var app = middleware.New(s.Install())

	app.Handle(newTestMux(s))
	server := httptest.NewServer(app)
	defer server.Close()
	jar, err = cookiejar.New(nil)
	if err != nil {
		panic(err)
	}
	client = &http.Client{}
	client.Jar = jar

	req, err = http.NewRequest("GET", server.URL+"/temporay", nil)
	if err != nil {
		panic(err)
	}
	resp, err = client.Do(req)
	if err != nil {
		panic(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Fatal(resp.StatusCode)
	}
	cookies = resp.Cookies()
	if len(cookies) != 1 || !cookies[0].Expires.IsZero() {
		t.Fatal(len(cookies), cookies)
	}

}
func TestNotAutoStart(t *testing.T) {
	var err error
	var s *Store
	var c *Cookie
	var client *http.Client
	var jar *cookiejar.Jar
	var req *http.Request
	var resp *http.Response
	s = newTestStore()
	c = &Cookie{
		httpcookie.Config{
			Name: "session",
		},
	}
	s.Installer = c
	var app = middleware.New(s.Install())

	app.Handle(newTestMux(s))
	server := httptest.NewServer(app)
	defer server.Close()
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
}

func catchPanic(h func()) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()
	h()
	return nil
}
