package httpcache

import (
	"net/http"

	"github.com/herb-go/datamodules/herbcache/cachepreset"
)

type PresetLoader interface {
	LoadPreset(r *http.Request) (*cachepreset.Preset, error)
}

type PresetLoaderFunc func(r *http.Request) (*cachepreset.Preset, error)

func (f PresetLoaderFunc) LoadPreset(r *http.Request) (*cachepreset.Preset, error) {
	return f(r)
}

func Preset(p *cachepreset.Preset) PresetLoader {
	return PresetLoaderFunc(func(r *http.Request) (*cachepreset.Preset, error) {
		return p, nil
	})
}
