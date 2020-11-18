package httpsession

import (
	"github.com/herb-go/herbdata"
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
	decrypted, err := AESNonceDecryptBase64(token, e.Secret)
	if err != nil {
		return "", nil, herbdata.ErrNotFound
	}
	return token, decrypted, nil
}
func (e *AESEngine) SaveToken(token string, data []byte, maxliftime int64) (newtoken string, err error) {
	encrypted, err := AESNonceEncryptBase64(data, e.Secret)
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
func EngineFactoryAES(loader func(v interface{}) error) (Engine, error) {
	e := &AESEngine{}
	err := loader(e)
	if err != nil {
		return nil, err
	}
	return e, nil
}
func init() {
	RegisterEngine(EngineNameAES, EngineFactoryAES)
}
