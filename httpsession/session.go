package httpsession

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
	Data       []*SessionValue
	CreatedAt  int64
	ExpiredAt  int64
	Temporay   bool
	LastActive int64
}

func newSessionData() *SessionData {
	return &SessionData{}
}

type Session struct {
	createdAt  *int64
	expiredAt  *int64
	lastactive *int64
	data       sync.Map
	updated    *int32
	token      atomic.Value
	started    *int32
	temporay   *int32
	loadedFrom string
}

func (s *Session) setdata(data *SessionData) {
	atomic.StoreInt64(s.createdAt, data.CreatedAt)
	atomic.StoreInt64(s.expiredAt, data.ExpiredAt)
	atomic.StoreInt64(s.lastactive, data.LastActive)
	for _, v := range data.Data {
		s.data.Store(v.Key, v.Value)
	}
	if data.Temporay {
		atomic.StoreInt32(s.temporay, 1)
	} else {
		atomic.StoreInt32(s.temporay, 0)
	}
}
func (s *Session) getdata() *SessionData {
	var data = newSessionData()
	data.CreatedAt = atomic.LoadInt64(s.createdAt)
	data.LastActive = atomic.LoadInt64(s.lastactive)
	data.ExpiredAt = atomic.LoadInt64(s.expiredAt)
	s.data.Range(func(key interface{}, val interface{}) bool {
		data.Data = append(data.Data, &SessionValue{Key: key.(string), Value: val.([]byte)})
		return true
	})
	data.Temporay = (atomic.LoadInt32(s.temporay) == 1)
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
func (s *Session) markAsNotUpdated() {
	atomic.StoreInt32(s.updated, 0)
}
func (s *Session) Updated() bool {
	updated := atomic.LoadInt32(s.updated)
	return updated == 1
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
func (s *Session) Temporay() bool {
	temporay := atomic.LoadInt32(s.temporay)
	return temporay == 1
}
func (s *Session) MarkAsTemporay() {
	s.MarkAsUpdated()
	atomic.StoreInt32(s.temporay, 1)
}

func (s *Session) LoadedFrom() string {
	return s.loadedFrom
}
func newSession() *Session {
	updated := int32(0)
	temporay := int32(0)
	started := int32(0)
	createdAt := int64(0)
	expiredAt := int64(0)
	lastactive := int64(0)

	return &Session{
		updated:    &updated,
		started:    &started,
		temporay:   &temporay,
		createdAt:  &createdAt,
		expiredAt:  &expiredAt,
		lastactive: &lastactive,
	}
}
