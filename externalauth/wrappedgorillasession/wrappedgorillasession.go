package wrappedgorillasession

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/sessions"
)

// Store session store
type Store struct {
	Store sessions.Store
	Name  string
}

// SessionRequestContextType content type which store session in request
type SessionRequestContextType string

// ContextField session request context type
var ContextField = SessionRequestContextType("Session")

// ErrNotFound error raised when session data not found
var ErrNotFound = errors.New("not found")

// New create new store with given  gorilla session and session context name
func New(store sessions.Store, name string) *Store {
	return &Store{
		Store: store,
		Name:  name,
	}
}

// Wrap wraper http handler with session
func (s *Store) Wrap(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cw := &CookieWriter{
			Writer: w.(Writer),
			r:      r,
			Store:  s,
		}

		s.MustLoad(cw, r)
		h.ServeHTTP(cw, r)
	})
}

// MustLoad load session from request.
//Panic if any error raised.
func (s *Store) MustLoad(w http.ResponseWriter, r *http.Request) {
	session, err := s.Store.Get(r, s.Name)
	if err != nil {
		panic(err)
	}
	s.SetSessionData(r, session)
}

// SetSessionData set session data to request.
func (s *Store) SetSessionData(r *http.Request, data *sessions.Session) {
	ctx := context.WithValue((*r).Context(), ContextField, data)
	req := (*r).WithContext(ctx)
	*r = *req
}

// SessionData get session data from request
func (s *Store) SessionData(r *http.Request) *sessions.Session {
	contextdata := r.Context().Value(ContextField)
	return contextdata.(*sessions.Session)
}

// Set set data to request with given field name
func (s *Store) Set(r *http.Request, fieldname string, v interface{}) error {
	bytes, err := json.Marshal(v)
	if err != nil {
		return err
	}
	data := s.SessionData(r)
	data.Values[fieldname] = bytes
	s.SetSessionData(r, data)
	return nil
}

//Get data from request with given field name.
//Param v must be pointer
func (s *Store) Get(r *http.Request, fieldname string, v interface{}) error {
	data := s.SessionData(r)
	value, ok := data.Values[fieldname]
	if ok == false {
		return ErrNotFound
	}
	bytes := value.([]byte)
	err := json.Unmarshal(bytes, v)
	return err
}

//Del delete data from request with given field name
func (s *Store) Del(r *http.Request, fieldname string) error {
	data := s.SessionData(r)
	delete(data.Values, fieldname)
	return nil
}

// IsNotFoundError check if give error is not found error
func (s *Store) IsNotFoundError(err error) bool {
	return err == ErrNotFound
}

// CookieWriter writer for cookie
type CookieWriter struct {
	Writer
	r       *http.Request
	Store   *Store
	written bool
}

// Writer http response writer interface
type Writer interface {
	http.ResponseWriter
	http.Hijacker
}

// WriteHeader write http response  header
func (w *CookieWriter) WriteHeader(status int) {
	if w.written == false {
		w.written = true
		w.Store.Store.Save(w.r, w, w.Store.SessionData(w.r))
	}
	w.Writer.WriteHeader(status)
}

// Write write http response  body
func (w *CookieWriter) Write(data []byte) (int, error) {
	if w.written == false {
		w.WriteHeader(http.StatusOK)
	}
	return w.Writer.Write(data)
}
