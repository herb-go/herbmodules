package httpsession

import (
	"github.com/herb-go/herbdata"
	"github.com/herb-go/herbsecurity/secret/encrypt/aesencrypt"
)

var EngineNameAES = EngineName("aes")

type AESEngine struct {
	Secret []byte
}

func (e *AESEngine) EngineName() EngineName {
	return EngineNameAES
}
func (e *AESEngine) NewToken() (token string, err error) {
	return TokenEmpty, nil
}
func (e *AESEngine) LoadToken(token string) (newtoken string, data []byte, err error) {
	if token == TokenEmpty {
		return "", nil, herbdata.ErrNotFound
	}
	decrypted, err := aesencrypt.AESNonceDecryptBase64(token, e.Secret)
	if err != nil {
		return "", nil, herbdata.ErrNotFound
	}
	return token, decrypted, nil
}
func (e *AESEngine) SaveToken(token string, data []byte, maxliftime int64) (newtoken string, err error) {
	encrypted, err := aesencrypt.AESNonceEncryptBase64(data, e.Secret)
	if err != nil {
		return "", err
	}
	return string(encrypted), nil
}
func (e *AESEngine) RevokeToken(token string) (err error) {
	return nil
}
func (e *AESEngine) DynamicToken() bool {
	return true
}
func (e *AESEngine) Start() error {
	return nil
}
func (e *AESEngine) Stop() error {
	return nil
}

type AESEngineConfig struct {
	Secret string
}

func (c *AESEngineConfig) CreateEngine() (Engine, error) {
	return &AESEngine{
		Secret: []byte(c.Secret),
	}, nil
}
func EngineFactoryAES(loader func(v interface{}) error) (Engine, error) {
	c := &AESEngineConfig{}
	err := loader(c)
	if err != nil {
		return nil, err
	}
	return c.CreateEngine()
}
func init() {
	RegisterEngine(EngineNameAES, EngineFactoryAES)
}
