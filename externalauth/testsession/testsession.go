package testsession

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"

	"golang.org/x/net/context"
)

var ErrNotFound = errors.New("not found")

type SessionRequestContextType string

var ContextField = SessionRequestContextType("Session")

type TestCookieSession struct {
	Name string
	Path string
}

type SessionData map[string][]byte

type CookieWriter struct {
	http.ResponseWriter
	r       *http.Request
	Session *TestCookieSession
	written bool
}

func (w *CookieWriter) WriteHeader(status int) {
	if w.written == false {
		w.written = true
		w.Session.MustSave(w, w.r)
	}
	w.ResponseWriter.WriteHeader(status)
}

func (w *CookieWriter) Write(data []byte) (int, error) {
	if w.written == false {
		w.WriteHeader(http.StatusOK)
	}
	return w.ResponseWriter.Write(data)
}

func (t *TestCookieSession) Wrap(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cw := &CookieWriter{
			ResponseWriter: w,
			r:              r,
			Session:        t,
		}

		t.MustLoad(cw, r)
		h.ServeHTTP(cw, r)
	})
}
func (t *TestCookieSession) MustLoad(w http.ResponseWriter, r *http.Request) {
	data := map[string][]byte{}
	cookie, err := r.Cookie(t.Name)
	if err == http.ErrNoCookie {
		t.SetSessionData(r, data)
		return
	}
	if err != nil {
		panic(err)
	}
	b64 := cookie.Value
	if b64 != "" {
		value, err := base64.StdEncoding.DecodeString(b64)
		if err != nil {
			panic(err)
		}
		err = json.Unmarshal(value, &data)
		if err != nil {
			panic(err)
		}
	}
	t.SetSessionData(r, data)
}

func (t *TestCookieSession) MustSave(w http.ResponseWriter, r *http.Request) {
	data := t.SessionData(r)
	bytes, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	b64 := base64.StdEncoding.EncodeToString(bytes)
	cookie := &http.Cookie{
		Name:  t.Name,
		Path:  t.Path,
		Value: b64,
	}
	http.SetCookie(w, cookie)
}

func (t *TestCookieSession) SetSessionData(r *http.Request, data SessionData) {
	ctx := context.WithValue((*r).Context(), ContextField, data)
	req := (*r).WithContext(ctx)
	*r = *req
}
func (t *TestCookieSession) SessionData(r *http.Request) SessionData {
	var data SessionData
	var ok bool
	contextdata := r.Context().Value(ContextField)
	if contextdata == nil {
		data = map[string][]byte{}
		return data
	} else if data, ok = contextdata.(SessionData); !ok {
		data = map[string][]byte{}
		return data
	}
	return data
}
func (t *TestCookieSession) Set(r *http.Request, fieldname string, v interface{}) error {
	bytes, err := json.Marshal(v)
	if err != nil {
		return err
	}
	data := t.SessionData(r)
	data[fieldname] = bytes
	t.SetSessionData(r, data)
	return nil
}
func (t *TestCookieSession) Get(r *http.Request, fieldname string, v interface{}) error {
	data := t.SessionData(r)
	bytes, ok := data[fieldname]
	if ok == false {
		return ErrNotFound
	}
	err := json.Unmarshal(bytes, v)
	return err
}
func (t *TestCookieSession) Del(r *http.Request, fieldname string) error {
	data := t.SessionData(r)
	delete(data, fieldname)
	return nil
}
func (t *TestCookieSession) IsNotFoundError(err error) bool {
	return err == ErrNotFound
}
