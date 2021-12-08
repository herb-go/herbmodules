package httpcache

import (
	"net/http"

	"github.com/herb-go/datamodules/herbcache"
	"github.com/herb-go/datamodules/herbcache/cachepreset"
	"github.com/herb-go/herb/identifier"
	"github.com/herb-go/herb/middleware/httpinfo"
	"github.com/herb-go/herbdata"
	"github.com/herb-go/herbdata/dataencoding/msgpackencoding"
)

var Encoding = msgpackencoding.Encoding

type Cache struct {
	Preset     *cachepreset.Preset
	Identifier identifier.Identifier
	Validator  httpinfo.Validator
}

func (c *Cache) Clone() *Cache {
	return &Cache{
		Preset:     c.Preset,
		Identifier: c.Identifier,
		Validator:  c.Validator,
	}
}
func (c *Cache) OverridePreset(p *cachepreset.Preset) *Cache {
	nc := c.Clone()
	nc.Preset = p
	return nc
}
func (c *Cache) OverrideIdentifier(i identifier.Identifier) *Cache {
	nc := c.Clone()
	nc.Identifier = i
	return nc
}
func (c *Cache) OverrideValidator(v httpinfo.Validator) *Cache {
	nc := c.Clone()
	nc.Validator = v
	return nc
}
func (c *Cache) ServeMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	if c.Preset == nil {
		next(w, r)
		return
	}
	key, err := c.Identifier.IdentifyRequest(r)
	if err != nil {
		panic(err)
	}
	if key == "" {
		next(w, r)
		return
	}
	var page = &cached{}
	resp := httpinfo.NewResponse()
	resp.UpdateAutocommit(false)
	var loader = func(key []byte) ([]byte, error) {
		controller := httpinfo.NewCommitController(r, resp).WithChecker(c.Validator)
		resp.UpdateController(controller)
		next(resp.WrapWriter(w), r)
		if resp.Autocommit() {
			return nil, herbcache.ErrNotCacheable
		}
		page = cacheResponse(resp)
		return Encoding.Marshal(page)
	}
	preset := c.Preset.Concat(cachepreset.Encoding(Encoding), cachepreset.Loader(loader))
	err = preset.Load([]byte(key), page)
	if err != nil {
		if err != herbdata.ErrEntryTooLarge && err != herbcache.ErrNotCacheable {
			panic(err)
		}
	}
	if !resp.Autocommit() {
		page.MustOutput(w)
	}
}

var DefaultValidator = httpinfo.ValidatorFunc(func(r *http.Request, resp *httpinfo.Response) (bool, error) {
	if !resp.Written {
		return true, nil
	}
	return resp.StatusCode >= 200 && resp.StatusCode < 500, nil
})

func New() *Cache {
	return &Cache{
		Validator:  DefaultValidator,
		Identifier: identifier.NopIdentifier,
	}
}
