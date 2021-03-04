package sharedrefresheraction

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type testRefresher struct {
	Data []byte
}

func (r *testRefresher) RefreshShared(old []byte) ([]byte, error) {
	if string(old) == "error" {
		return nil, errors.New("noperr")
	}
	return r.Data, nil
}
func TestAction(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if r := recover(); r != nil {
					http.Error(w, http.StatusText(500), 500)
				}
			}()
			NewSharedRefresherAction(&testRefresher{Data: []byte("testdata")}).ServeHTTP(w, r)
		}),
	)
	defer s.Close()
	req, err := http.NewRequest("GET", s.URL, nil)
	if err != nil {
		panic(err)
	}
	resp, err := s.Client().Do(req)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != 405 {
		t.Fatal()
	}
	resp.Body.Close()
	req, err = http.NewRequest("POST", s.URL, bytes.NewBufferString("test"))
	if err != nil {
		panic(err)
	}
	resp, err = s.Client().Do(req)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode != 200 {
		t.Fatal()
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	resp.Body.Close()
	if bytes.Compare(data, []byte("testdata")) != 0 {
		t.Fatal()
	}
	req, err = http.NewRequest("POST", s.URL, bytes.NewBufferString("error"))
	if err != nil {
		panic(err)
	}
	resp, err = s.Client().Do(req)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode != 500 {
		t.Fatal()
	}
	resp.Body.Close()
}
