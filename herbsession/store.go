package herbsession

import (
	"context"
	"net/http"
	"time"

	"github.com/herb-go/herbdata"
	"gopkg.in/vmihailenco/msgpack.v2"
)

type SessionName string

type Store struct {
	Name        SessionName
	AutoStart   bool
	engine      Engine
	MaxLifetime int64
}

func (s *Store) StartSession() (string, *Session, error) {
	t, err := s.engine.NewToken()
	if err != nil {
		return "", nil, err
	}
	session := newSession()
	session.MarkAsStarted()
	session.token.Store(t)
	now := time.Now().Unix()
	session.expiredAt = now + s.MaxLifetime
	session.createdAt = now
	return t, session, nil
}

func (s *Store) LoadSession(token string) (*Session, error) {
	t, data, err := s.engine.LoadToken(token)
	if err != nil {
		return nil, err
	}
	sessiondata := newSessionData()
	err = msgpack.Unmarshal(data, sessiondata)
	if err != nil {
		return nil, err
	}
	if sessiondata.ExpiredAt < time.Now().Unix() {
		return nil, herbdata.ErrNotFound
	}
	session := newSession()
	session.MarkAsStarted()
	session.setdata(sessiondata)
	session.token.Store(t)
	if t != token {
		session.MarkAsUpdated()
	}
	return session, nil
}

func (s *Store) SaveSession(session *Session) (newtoken string, err error) {
	data, err := msgpack.Marshal(session.getdata())
	if err != nil {
		return "", err
	}
	return s.engine.UpdateToken(session.Token(), data, s.MaxLifetime)
}
func (s *Store) RevokeSession(token string) (newtoken string, err error) {
	return s.engine.RevokeToken(token)
}

func (s *Store) SessionLastActive(token string) (int64, error) {
	return s.engine.TokenLastActive(token)
}

func (s *Store) RequestSession(r *http.Request) (session *Session) {
	v := r.Context().Value(s.Name)
	if v == nil {
		return nil
	}
	return v.(*Session)
}
func (s *Store) SetRequestSession(r *http.Request, session *Session) {
	ctx := context.WithValue(r.Context(), s.Name, session)
	*r = *(r.WithContext(ctx))
}
