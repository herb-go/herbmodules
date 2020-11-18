package httpsession

import (
	"context"
	"net/http"
	"time"

	"github.com/herb-go/herbdata"
	"gopkg.in/vmihailenco/msgpack.v2"
)

type SessionName string

type Store struct {
	Name               SessionName
	AutoStart          bool
	Engine             Engine
	MaxLifetime        int64
	Timeout            int64
	LastActiveInterval int64
	Installer          Installer
}

func (s *Store) Install() func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	return s.Installer.InstallerMiddleware(s)
}

func (s *Store) StartSession() (string, *Session, error) {
	t, err := s.Engine.NewToken()
	if err != nil {
		return "", nil, err
	}
	session := newSession()
	session.MarkAsStarted()
	session.MarkAsUpdated()
	session.token.Store(t)
	now := time.Now().Unix()
	expiredAt := now + s.MaxLifetime
	createdAt := now
	lastactive := now
	session.expiredAt = &expiredAt
	session.createdAt = &createdAt
	session.lastactive = &lastactive
	return t, session, nil
}

func (s *Store) LoadSession(token string) (*Session, error) {
	t, data, err := s.Engine.LoadToken(token)
	if err != nil {
		return nil, err
	}
	sessiondata := newSessionData()
	err = msgpack.Unmarshal(data, sessiondata)
	if err != nil {
		return nil, err
	}
	now := time.Now().Unix()
	if sessiondata.ExpiredAt < now {
		return nil, herbdata.ErrNotFound
	}
	if s.Timeout > 0 && sessiondata.LastActive+s.Timeout < now {
		return nil, herbdata.ErrNotFound
	}
	session := newSession()
	session.MarkAsStarted()
	if s.Timeout > 0 && sessiondata.LastActive+s.LastActiveInterval < now {
		sessiondata.LastActive = now
		session.MarkAsUpdated()
	}
	session.setdata(sessiondata)
	session.token.Store(t)
	session.loadedFrom = token
	if t != token {
		session.MarkAsUpdated()
	}
	return session, nil
}

func (s *Store) SaveSession(session *Session) (err error) {
	data, err := msgpack.Marshal(session.getdata())
	if err != nil {
		return err
	}
	newtoken, err := s.Engine.SaveToken(session.Token(), data, s.MaxLifetime)
	if err != nil {
		return err
	}
	session.SetToken(newtoken)
	return nil
}
func (s *Store) RevokeSession(token string) (err error) {
	return s.Engine.RevokeToken(token)
}

func (s *Store) RequestSession(r *http.Request) (session *Session) {
	v := r.Context().Value(s.Name)
	if v == nil {
		return nil
	}
	return v.(*Session)
}
func (s *Store) RegenerateRequestSessionID(r *http.Request) error {
	token, err := s.Engine.NewToken()
	if err != nil {
		return err
	}
	session := s.RequestSession(r)
	session.SetToken(token)
	session.MarkAsUpdated()
	return nil
}
func (s *Store) MustInstallSessionToRequest(rpointer **http.Request, token string) *Session {
	var err error
	var session *Session
	if token != "" {
		session, err = s.LoadSession(token)
		if err != nil {
			panic(err)
		}
	}
	if session == nil {
		if s.AutoStart {
			_, session, err = s.StartSession()
		} else {
			session = newSession()
		}
	}
	s.SetRequestSession(rpointer, session)
	return session
}
func (s *Store) RevokeRequestSession(r *http.Request) (err error) {
	session := s.RequestSession(r)
	token := session.Token()
	if token == TokenEmpty {
		return nil
	}
	err = s.Engine.RevokeToken(token)
	if err != nil {
		return err
	}
	session.token.Store(TokenEmpty)
	session.markAsNotUpdated()
	return nil

}
func (s *Store) SetRequestSession(r **http.Request, session *Session) {
	ctx := context.WithValue((*r).Context(), s.Name, session)
	req := (*r).WithContext(ctx)
	*r = req
}
func (s *Store) Start() error {
	return s.Engine.Start()
}
func (s *Store) Stop() error {
	return s.Engine.Stop()
}
func New() *Store {
	return &Store{}
}
