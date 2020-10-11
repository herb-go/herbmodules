package commonextractor

import (
	"net"
	"net/http"
	"net/url"

	"github.com/herb-go/herb/middleware/httpinfo"
	"github.com/herb-go/herb/middleware/router"
	"github.com/herb-go/herbmodules/httpinfomanager"
)

type FieldConfig struct {
	Field string
}

type HeaderExtractor string

func (e HeaderExtractor) Extract(r *http.Request) ([]byte, error) {
	return []byte(r.Header.Get(string(e))), nil
}

var HeaderExtractorFactory = httpinfomanager.ExtractorFactoryFunc(func(loader func(interface{}) error) (httpinfo.Extractor, error) {
	c := &FieldConfig{}
	err := loader(c)
	if err != nil {
		return nil, err
	}
	return HeaderExtractor(c.Field), nil
})

type QueryExtractor string

func (e QueryExtractor) Extract(r *http.Request) ([]byte, error) {
	q := r.URL.Query()
	return []byte(q.Get(string(e))), nil
}

var QueryExtractorFactory = httpinfomanager.ExtractorFactoryFunc(func(loader func(interface{}) error) (httpinfo.Extractor, error) {
	c := &FieldConfig{}
	err := loader(c)
	if err != nil {
		return nil, err
	}
	return QueryExtractor(c.Field), nil
})

type FormExtractor string

func (e FormExtractor) Extract(r *http.Request) ([]byte, error) {
	v := r.FormValue(string(e))
	return []byte(v), nil
}

var FormExtractorFactory = httpinfomanager.ExtractorFactoryFunc(func(loader func(interface{}) error) (httpinfo.Extractor, error) {
	c := &FieldConfig{}
	err := loader(c)
	if err != nil {
		return nil, err
	}
	return FormExtractor(c.Field), nil
})

type RouterExtractor string

func (e RouterExtractor) Extract(r *http.Request) ([]byte, error) {
	p := router.GetParams(r)
	return []byte(p.Get(string(e))), nil
}

var RouterExtractorFactory = httpinfomanager.ExtractorFactoryFunc(func(loader func(interface{}) error) (httpinfo.Extractor, error) {
	c := &FieldConfig{}
	err := loader(c)
	if err != nil {
		return nil, err
	}
	return RouterExtractor(c.Field), nil
})

type FixedExtractor string

func (e FixedExtractor) Extract(r *http.Request) ([]byte, error) {
	return []byte(string(e)), nil
}

var FixedExtractorFactory = httpinfomanager.ExtractorFactoryFunc(func(loader func(interface{}) error) (httpinfo.Extractor, error) {
	c := &FieldConfig{}
	err := loader(c)
	if err != nil {
		return nil, err
	}
	return FixedExtractor(c.Field), nil
})

type CookieExtractor string

func (e CookieExtractor) Extract(r *http.Request) ([]byte, error) {
	c, err := r.Cookie(string(e))
	if err != nil {
		if err == http.ErrNoCookie {
			return nil, nil
		}
		return nil, err
	}
	return []byte(c.Value), nil
}

var CookieExtractorFactory = httpinfomanager.ExtractorFactoryFunc(func(loader func(interface{}) error) (httpinfo.Extractor, error) {
	c := &FieldConfig{}
	err := loader(c)
	if err != nil {
		return nil, err
	}
	return CookieExtractor(c.Field), nil
})

var IPAddressExtractor = httpinfo.ExtractorFunc(func(r *http.Request) ([]byte, error) {
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return nil, err
	}
	return []byte(ip), nil
})

var IPAddressExtractorFactory = httpinfomanager.ExtractorFactoryFunc(func(loader func(v interface{}) error) (httpinfo.Extractor, error) {
	return IPAddressExtractor, nil
})

var MethodExtractor = httpinfo.ExtractorFunc(func(r *http.Request) ([]byte, error) {
	return []byte(r.Method), nil
})

var MethodExtractorFactory = httpinfomanager.ExtractorFactoryFunc(func(loader func(v interface{}) error) (httpinfo.Extractor, error) {
	return MethodExtractor, nil
})

var PathExtractor = httpinfo.ExtractorFunc(func(r *http.Request) ([]byte, error) {
	p := r.RequestURI
	u, err := url.Parse(p)
	if err != nil {
		return nil, err
	}
	return []byte(u.Path), nil
})

var PathExtractorFactory = httpinfomanager.ExtractorFactoryFunc(func(loader func(v interface{}) error) (httpinfo.Extractor, error) {
	return PathExtractor, nil
})

var HostExtractor = httpinfo.ExtractorFunc(func(r *http.Request) ([]byte, error) {
	return []byte(r.Host), nil
})

var HostExtractorFactory = httpinfomanager.ExtractorFactoryFunc(func(loader func(v interface{}) error) (httpinfo.Extractor, error) {
	return HostExtractor, nil
})

var UserExtractor = httpinfo.ExtractorFunc(func(r *http.Request) ([]byte, error) {
	u, _, _ := r.BasicAuth()
	return []byte(u), nil
})

var UserExtractorFactory = httpinfomanager.ExtractorFactoryFunc(func(loader func(v interface{}) error) (httpinfo.Extractor, error) {
	return UserExtractor, nil
})

var PasswordExtractor = httpinfo.ExtractorFunc(func(r *http.Request) ([]byte, error) {
	_, p, _ := r.BasicAuth()
	return []byte(p), nil
})

var PasswordExtractorFactory = httpinfomanager.ExtractorFactoryFunc(func(loader func(v interface{}) error) (httpinfo.Extractor, error) {
	return PasswordExtractor, nil
})

func RegsiterFactories() {
	httpinfomanager.RegisterExtractorFactory("header", HeaderExtractorFactory)
	httpinfomanager.RegisterExtractorFactory("query", QueryExtractorFactory)
	httpinfomanager.RegisterExtractorFactory("form", FormExtractorFactory)
	httpinfomanager.RegisterExtractorFactory("router", RouterExtractorFactory)
	httpinfomanager.RegisterExtractorFactory("fixed", FixedExtractorFactory)
	httpinfomanager.RegisterExtractorFactory("cookie", CookieExtractorFactory)
	httpinfomanager.RegisterExtractorFactory("ip", IPAddressExtractorFactory)
	httpinfomanager.RegisterExtractorFactory("method", MethodExtractorFactory)
	httpinfomanager.RegisterExtractorFactory("path", PathExtractorFactory)
	httpinfomanager.RegisterExtractorFactory("host", HostExtractorFactory)
	httpinfomanager.RegisterExtractorFactory("user", UserExtractorFactory)
	httpinfomanager.RegisterExtractorFactory("password", PasswordExtractorFactory)
}
