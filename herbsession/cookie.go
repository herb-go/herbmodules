package herbsession

import (
	"net/http"
	"time"

	"github.com/herb-go/herb/middleware"
	"github.com/herb-go/herb/service/httpservice/httpcookie"
)

var InstallerIDCookie = InstallerID("cookie")

type Cookie struct {
	Config httpcookie.Config
}

func (c *Cookie) InstallerMiddleware(s *Store) func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		var token string
		cookie, err := r.Cookie(c.Config.Name)
		if err == nil {
			token = cookie.Value
		} else if err != http.ErrNoCookie {
			panic(err)
		}
		session := s.MustInstallSessionToRequest(&r, token)
		writer := &cookiewriter{
			store:          s,
			session:        session,
			cookie:         c,
			ResponseWriter: w,
		}
		next(middleware.WrapResponseWriter(writer), r)
	}
}

type cookiewriter struct {
	store   *Store
	session *Session
	http.ResponseWriter
	cookie  *Cookie
	written bool
}

func (w *cookiewriter) WriteHeader(status int) {
	var err error
	if w.written == false {
		w.written = true
		if w.session.Updated() {

			err = w.store.SaveSession(w.session)
			if err != nil {
				panic(err)
			}
		}
		token := w.session.Token()
		if token != w.session.LoadedFrom() {
			cookie := w.cookie.Config.CreateCookieWithValue(token)
			if token != "" {
				if !w.session.Temporay() {
					cookie.Expires = time.Now().Add(time.Duration(w.store.MaxLifetime) * time.Second)
				}
			} else {
				cookie.Expires = time.Unix(0, 0)
			}
			http.SetCookie(w, cookie)
		}
	}
	w.ResponseWriter.WriteHeader(status)
}

func (w *cookiewriter) Write(data []byte) (int, error) {
	if w.written == false {
		w.WriteHeader(200)
	}
	return w.ResponseWriter.Write(data)
}
