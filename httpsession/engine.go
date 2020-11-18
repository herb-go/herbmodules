package httpsession

import "sync"

type EngineName string
type Engine interface {
	EngineName() EngineName
	NewToken() (token string, err error)
	LoadToken(token string) (newtoken string, data []byte, err error)
	SaveToken(token string, data []byte, maxLiftimeInSecond int64) (newtoken string, err error)
	RevokeToken(token string) (err error)
	DynamicToken() bool
	Start() error
	Stop() error
}

type EngineFactory func(func(v interface{}) error) (Engine, error)

var enginesFactories sync.Map

func RegisterEngine(name EngineName, f EngineFactory) {
	enginesFactories.Store(name, f)
}
func CreateEngine(name EngineName, loader func(v interface{}) error) (Engine, error) {
	v, ok := enginesFactories.Load(name)
	if !ok {
		return nil, NewUnknownSessionEngineError(name)
	}
	return v.(EngineFactory)(loader)
}
