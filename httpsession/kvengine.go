package httpsession

import (
	"crypto/rand"
	"encoding/base64"
	"time"

	"github.com/herb-go/herbdata/kvdb"
	"gopkg.in/vmihailenco/msgpack.v2"
)

const kvPrefixData = "\x00"
const kvPrefixLastActive = "\x01"

var EngineNameKV = EngineName("keyvalue")

type KVEngine struct {
	Database  *kvdb.Database
	TokenSize int
	Timeout   int64
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
func (e *KVEngine) SessionTimeout() int64 {
	if e.Timeout <= 0 {
		return 0
	}
	return e.Timeout
}
func (e *KVEngine) TokenLastActive(token string) (int64, error) {

	data, err := e.Database.Get([]byte(kvPrefixLastActive + token))
	if err != nil {
		return 0, err
	}
	var lastactive int64
	err = msgpack.Unmarshal(data, &lastactive)
	if err != nil {
		return 0, err
	}
	return lastactive, nil
}
func (e *KVEngine) LoadToken(token string) (newtoken string, data []byte, err error) {
	data, err = e.Database.Get([]byte(kvPrefixData + token))
	if err != nil {
		return "", nil, err
	}
	if e.Timeout > 0 {
		lastactivedata, err := msgpack.Marshal(time.Now().Unix())
		if err != nil {
			return "", nil, err
		}
		err = e.Database.SetWithTTL([]byte(kvPrefixLastActive+token), lastactivedata, e.Timeout)
		if err != nil {
			return "", nil, err
		}
	}
	return token, data, nil
}
func (e *KVEngine) SaveToken(token string, data []byte, maxLiftimeInSecond int64) (newtoken string, err error) {
	err = e.Database.SetWithTTL([]byte(kvPrefixData+token), data, maxLiftimeInSecond)
	if err != nil {
		return "", err
	}
	if e.Timeout > 0 {
		lastactivedata, err := msgpack.Marshal(time.Now().Unix())
		if err != nil {
			return "", err
		}
		err = e.Database.SetWithTTL([]byte(kvPrefixLastActive+token), lastactivedata, e.Timeout)
		if err != nil {
			return "", err
		}
	}
	return token, err
}
func (e *KVEngine) RevokeToken(token string) (err error) {
	return e.Database.Delete([]byte(kvPrefixData + token))
}
func (e *KVEngine) DynamicToken() bool {
	return false
}
func (e *KVEngine) Close() error {
	return e.Database.Close()
}

type KVEngineConfig struct {
	kvdb.Config
	TokenSize int
	Timeout   int64
}

func (c *KVEngineConfig) ApplyTo(e *KVEngine) error {
	db := kvdb.New()
	err := c.Config.ApplyTo(db)
	if err != nil {
		return err
	}
	e.Database = db
	e.TokenSize = c.TokenSize
	e.Timeout = c.Timeout
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
