package dchest

import (
	"time"

	"github.com/herb-go/herbmodules/captcha/drivers/dchest/dchestimage"

	"github.com/herb-go/herbmodules/captcha"
	"github.com/herb-go/herb/cache"
)

type Config struct {
	Secret      string
	Cache       *cache.OptionConfig
	CachePrefix string
	Size        int
	TTLInSecond int64
	Width       int
	Height      int
}

func (c *Config) Create() (*Driver, error) {
	var err error
	ttl := DefaultTTL
	if c.TTLInSecond != 0 {
		ttl = time.Duration(c.TTLInSecond) * time.Second
	}
	c.Cache.TTL = int64(ttl / time.Second)
	tokencache := cache.New()
	err = c.Cache.ApplyTo(tokencache)
	if err != nil {
		return nil, err
	}
	d := &Driver{
		Cache:       tokencache,
		CachePrefix: c.CachePrefix,
		Secret:      c.Secret,
		Wanted: &captcha.Wanted{
			Min:          DefaultSize,
			Max:          DefaultSize,
			OptionalByte: DefaultOptionalBytes,
		},
		ImageRender:    dchestimage.DefaultConfig,
		TTL:            ttl,
		Signer:         DefaultSigner,
		TokenGenerator: DefaultTokenGenerator,
	}
	if c.Width == 0 {
		d.Width = DefaultWidth
	}
	if c.Height == 0 {
		d.Height = DefaultHeight
	}
	return d, nil
}
