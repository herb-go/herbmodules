package wrappedgorillasession

import (
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/sessions"
)

func TestSession(t *testing.T) {
	gstore := sessions.NewCookieStore([]byte("something-very-secret"))
	session := New(gstore, "test")

	mux := http.NewServeMux()
	mux.Handle("/set", session.Wrap(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			value := r.URL.Query().Get("value")
			err := session.Set(r, "test", value)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			w.Write([]byte("ok"))
		},
	)))
	mux.Handle("/del", session.Wrap(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			err := session.Del(r, "test")
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			w.Write([]byte("ok"))

			w.Write([]byte("ok"))
		},
	)))
	mux.Handle("/get", session.Wrap(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			var value string
			err := session.Get(r, "test", &value)
			if session.IsNotFoundError(err) {
				http.Error(w, http.StatusText(404), 404)
				return
			}
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			w.Write([]byte(value))
		},
	)))
	server := httptest.NewServer(mux)
	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Error(err)
	}
	client := http.Client{
		Jar: jar,
	}
	req, err := http.NewRequest("post", server.URL+"/get", nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 404 {
		t.Error(resp.StatusCode)
	}
	req, err = http.NewRequest("post", server.URL+"/set?value="+"test", nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err = client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Error(resp.StatusCode)
	}
	req, err = http.NewRequest("post", server.URL+"/get", nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err = client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Error(resp.StatusCode)
	}
	if string(content) != "test" {
		t.Error(string(content))
	}

	req, err = http.NewRequest("post", server.URL+"/set?value="+"test2", nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err = client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Error(resp.StatusCode)
	}
	req, err = http.NewRequest("post", server.URL+"/get", nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err = client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	content, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Error(resp.StatusCode)
	}
	if string(content) != "test2" {
		t.Error(string(content))
	}

	req, err = http.NewRequest("post", server.URL+"/del", nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err = client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Error(resp.StatusCode)
	}

	req, err = http.NewRequest("post", server.URL+"/get", nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err = client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 404 {
		t.Error(resp.StatusCode)
	}
}
