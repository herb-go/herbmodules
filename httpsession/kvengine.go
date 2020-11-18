package httpsession

import (
	"crypto/rand"
	"encoding/base64"

	"github.com/herb-go/herbdata/kvdb"
)

var EngineNameKV = EngineName("keyvalue")

type KVEngine struct {
	Database  *kvdb.Database
	TokenSize int
}

func (e *KVEngine) EngineName() EngineName {
	return EngineNameKV
}
func (e *KVEngine) NewToken() (token string, err error) {
	var tokendata = make([]byte, e.TokenSize)
	_, err = rand.Read(tokendata)
	if err != nil {
		return "", nil
	}
	return base64.StdEncoding.EncodeToString(tokendata), nil
}

func (e *KVEngine) LoadToken(token string) (newtoken string, data []byte, err error) {
	data, err = e.Database.Get([]byte(token))
	if err != nil {
		return "", nil, err
	}
	return token, data, nil
}
func (e *KVEngine) SaveToken(token string, data []byte, maxLifetimeInSecond int64) (newtoken string, err error) {
	err = e.Database.SetWithTTL([]byte(token), data, maxLifetimeInSecond)
	if err != nil {
		return "", err
	}
	return token, err
}
func (e *KVEngine) RevokeToken(token string) (err error) {
	return e.Database.Delete([]byte(token))
}
func (e *KVEngine) DynamicToken() bool {
	return false
}
func (e *KVEngine) Start() error {
	return e.Database.Start()
}
func (e *KVEngine) Stop() error {
	return e.Database.Stop()
}

type KVEngineConfig struct {
	kvdb.Config
	TokenSize int
}

func (c *KVEngineConfig) ApplyTo(e *KVEngine) error {
	db := kvdb.New()
	err := c.Config.ApplyTo(db)
	if err != nil {
		return err
	}
	e.Database = db
	e.TokenSize = c.TokenSize
	return nil
}

func EngineFactoryKV(loader func(v interface{}) error) (Engine, error) {
	e := &KVEngine{}
	c := &KVEngineConfig{}
	err := loader(c)
	if err != nil {
		return nil, err
	}
	err = c.ApplyTo(e)
	if err != nil {
		return nil, err
	}
	return e, nil
}
func init() {
	RegisterEngine(EngineNameKV, EngineFactoryKV)
}
