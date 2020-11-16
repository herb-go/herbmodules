package herbsession

import (
	"sync"
	"sync/atomic"

	"github.com/herb-go/herbdata"
)

type SessionValue struct {
	Key   string
	Value []byte
}

type SessionData struct {
	Data      []*SessionValue
	CreatedAt int64
	ExpiredAt int64
}

func newSessionData() *SessionData {
	return &SessionData{}
}

type Session struct {
	createdAt int64
	expiredAt int64
	data      sync.Map
	updated   *int32
	revoked   *int32
	token     atomic.Value
	started   *int32
}

func (s *Session) setdata(data *SessionData) {
	s.createdAt = data.CreatedAt
	s.expiredAt = data.ExpiredAt
	for _, v := range data.Data {
		s.data.Store(v.Key, v.Value)
	}
}
func (s *Session) getdata() *SessionData {
	var data = newSessionData()
	data.CreatedAt = s.createdAt
	data.ExpiredAt = s.expiredAt
	s.data.Range(func(key interface{}, val interface{}) bool {
		data.Data = append(data.Data, &SessionValue{Key: key.(string), Value: val.([]byte)})
		return true
	})
	return data
}

func (s *Session) Store(key string, data []byte) error {
	if !s.Started() {
		return ErrSessionNotStarted
	}
	s.data.Store(key, data)
	s.MarkAsUpdated()
	return nil
}

func (s *Session) Set(key []byte, data []byte) error {
	return s.Store(string(key), data)
}
func (s *Session) Load(key string) ([]byte, error) {
	if !s.Started() {
		return nil, ErrSessionNotStarted
	}
	v, ok := s.data.Load(key)
	if !ok {
		return nil, herbdata.ErrNotFound
	}
	return v.([]byte), nil

}
func (s *Session) Get(key []byte) ([]byte, error) {
	return s.Load(string(key))
}

func (s *Session) Remove(key string) error {
	if !s.Started() {
		return ErrSessionNotStarted
	}
	s.data.Delete(key)
	s.MarkAsUpdated()
	return nil
}
func (s *Session) Delete(key []byte) error {
	return s.Remove(string(key))
}
func (s *Session) MarkAsUpdated() {
	atomic.StoreInt32(s.updated, 1)
}
func (s *Session) Updated() bool {
	updated := atomic.LoadInt32(s.updated)
	return updated == 1
}
func (s *Session) MarkAsRevoked() {
	atomic.StoreInt32(s.revoked, 1)
	s.MarkAsUpdated()
}
func (s *Session) Revoked() bool {
	revoked := atomic.LoadInt32(s.revoked)
	return revoked == 1
}
func (s *Session) MarkAsStarted() {
	atomic.StoreInt32(s.started, 1)
}
func (s *Session) Started() bool {
	started := atomic.LoadInt32(s.started)
	return started == 1
}
func (s *Session) SetToken(token string) {
	s.token.Store(token)
}
func (s *Session) Token() string {
	var token string
	t := s.token.Load()
	if t != nil {
		token = t.(string)
	}
	return token
}
func newSession() *Session {
	updated := int32(0)
	revoked := int32(0)
	started := int32(0)
	return &Session{
		updated: &updated,
		revoked: &revoked,
		started: &started,
	}
}
