package httpcache

import (
	"net/http"

	"github.com/herb-go/herb/middleware/httpinfo"
)

type cached struct {
	StatusCode int
	Header     http.Header
	Body       []byte
}

func mergeHeader(src http.Header, dst *http.Header) {
	for key := range src {
		for k := range src[key] {
			(*dst).Add(key, src[key][k])
		}
	}
}

func (c *cached) MustOutput(w http.ResponseWriter) {
	h := w.Header()
	mergeHeader(c.Header, &h)
	w.WriteHeader(c.StatusCode)
	_, err := w.Write(c.Body)
	if err != nil {
		panic(err)
	}
}

func cacheResponse(resp *httpinfo.Response) *cached {
	data := resp.UncommittedData()
	c := &cached{
		StatusCode: resp.StatusCode,
		Header:     http.Header{},
		Body:       data,
	}
	mergeHeader(resp.Header(), &c.Header)
	return c
}
