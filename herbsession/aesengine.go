package herbsession

import (
	"time"

	"github.com/herb-go/herbdata"
	"gopkg.in/vmihailenco/msgpack.v2"
)

type AESEngine struct {
	Secret  []byte
	Timeout int64
}

type aesEngineData struct {
	LastActive int64
	Data       []byte
}

func (e *AESEngine) NewToken() (token string, err error) {
	if e.Timeout <= 0 {
		return TokenEmpty, nil
	}
	d := &aesEngineData{
		LastActive: time.Now().Unix(),
	}
	bytes, err := msgpack.Marshal(d)
	if err != nil {
		return "", err
	}

	return AESNonceEncryptBase64(bytes, e.Secret)
}
func (e *AESEngine) TokenLastActive(token string) (int64, error) {
	if e.Timeout <= 0 {
		return 0, nil
	}
	if token == TokenEmpty {
		return 0, herbdata.ErrNotFound
	}
	data, err := AESNonceDecryptBase64(token, e.Secret)
	if err != nil {
		return 0, herbdata.ErrNotFound
	}
	d := &aesEngineData{}
	err = msgpack.Unmarshal(data, d)
	if err != nil {
		return 0, err
	}
	if d.LastActive+e.Timeout < time.Now().Unix() {
		return 0, herbdata.ErrNotFound
	}
	return d.LastActive, nil
}
func (e *AESEngine) LoadToken(token string) (newtoken string, data []byte, err error) {
	if token == TokenEmpty {
		return "", nil, herbdata.ErrNotFound
	}
	decrypted, err := AESNonceDecryptBase64(token, e.Secret)
	if err != nil {
		return "", nil, herbdata.ErrNotFound
	}
	d := &aesEngineData{}
	err = msgpack.Unmarshal(decrypted, d)
	if err != nil {
		return "", nil, herbdata.ErrNotFound
	}
	if e.Timeout <= 0 {
		return token, d.Data, nil
	}
	if d.LastActive+e.Timeout < time.Now().Unix() {
		return "", nil, herbdata.ErrNotFound
	}
	newdata := aesEngineData{
		LastActive: time.Now().Unix(),
		Data:       d.Data,
	}
	bytes, err := msgpack.Marshal(newdata)
	if err != nil {
		return "", nil, err
	}
	nt, err := AESNonceEncryptBase64(bytes, e.Secret)
	if err != nil {
		return "", nil, err
	}
	return nt, d.Data, nil
}
func (e *AESEngine) UpdateToken(token string, data []byte, maxliftime int64) (newtoken string, err error) {
	aesdata := aesEngineData{
		LastActive: time.Now().Unix(),
		Data:       data,
	}
	bytes, err := msgpack.Marshal(aesdata)
	if err != nil {
		return "", err
	}
	encrypted, err := AESNonceEncryptBase64(bytes, e.Secret)
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
func (e *AESEngine) SessionTimeout() int64 {
	return e.Timeout
}
