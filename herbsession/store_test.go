package herbsession

import (
	"bytes"
	"testing"
	"time"

	"github.com/herb-go/herbdata"
)

func TestStore(t *testing.T) {
	var s *Store
	var err error
	var session *Session
	var newsession *Session
	var data []byte
	var token string
	var lastactive int64
	var e = newTestAESEngine()
	e.Timeout = 100
	s = &Store{
		Name:        "test",
		MaxLifetime: 2,
		AutoStart:   false,
		engine:      e,
	}

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
	token, err = s.SaveSession(session)
	if err != nil {
		t.Fatal()
	}
	time.Sleep(1 * time.Second)
	newsession, err = s.LoadSession(token)
	if err != nil {
		t.Fatal()
	}
	data, err = newsession.Load("test")
	if err != nil {
		t.Fatal()
	}
	if bytes.Compare(data, []byte("testdata")) != 0 {
		t.Fatal()
	}
	lastactive, err = s.SessionLastActive(newsession.Token())

	if err != nil || lastactive-time.Now().Unix() > 1 {
		t.Fatal()
	}
	token, err = s.RevokeSession(newsession.Token())
	if err != nil {
		t.Fatal()
	}
	newsession, err = s.LoadSession(token)
	if err != herbdata.ErrNotFound {
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
	token, err = s.SaveSession(session)
	if err != nil {
		t.Fatal()
	}
	_, err = s.LoadSession(token)
	if err != nil {
		t.Fatal()
	}
	time.Sleep(3 * time.Second)
	_, err = s.LoadSession(token)
	if err != herbdata.ErrNotFound {
		t.Fatal()
	}
}
