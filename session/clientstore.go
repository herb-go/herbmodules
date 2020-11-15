package session

import (
	"bytes"
	"crypto/rand"
	"time"

	"github.com/herb-go/deprecated/cache"

	"crypto/aes"
	"crypto/cipher"

	"encoding/base64"
)

const clientStoreNonceSize = 4
const clientStoreNewToken = "."

var filledByte = []byte{0}

//IVSize AES IV size
const IVSize = 16

func formatKey(key []byte, size int) []byte {
	var data = make([]byte, size)
	copy(data, key)
	return data
}

//AESEncrypt aes encrypt with given data,key and iv.
//Data will be padding with PKCS7Padding
//Return encrytped data and any error if raised.
func AESEncrypt(unencrypted []byte, key []byte, iv []byte) (encrypted []byte, err error) {
	defer func() {
		r := recover()
		if r != nil {
			err = r.(error)
		}
	}()
	cryptKey := formatKey(key, aes.BlockSize)
	block, err := aes.NewCipher(cryptKey)
	if err != nil {
		return
	}
	data := PKCS7Padding(unencrypted, aes.BlockSize)
	crypter := cipher.NewCBCEncrypter(block, iv)
	encrypted = make([]byte, len(data))
	crypter.CryptBlocks(encrypted, data)
	return
}

// AESNonceEncrypt aes encrypt data with given key and random bytes as IV.
//Data will be padding with PKCS7Padding
//Random IV will prefix encryped data
//return encrypted data and any error if raisd.
func AESNonceEncrypt(unencrypted []byte, key []byte) (encrypted []byte, err error) {
	defer func() {
		r := recover()
		if r != nil {
			err = r.(error)
		}
	}()
	var rawEncrypted []byte
	var IV = make([]byte, IVSize)
	_, err = rand.Read(IV)
	if err != nil {
		return
	}
	rawEncrypted, err = AESEncrypt(unencrypted, key, IV)
	if err != nil {
		return
	}
	encrypted = make([]byte, len(rawEncrypted)+int(IVSize))
	copy(encrypted[:IVSize], IV)
	copy(encrypted[IVSize:], rawEncrypted)
	return
}

//AESEncryptBase64 aes encrypt with given data,key and iv.
//Data will be padding with PKCS7Padding
//Return base64 encoded encrytped data and any error if raised.
func AESEncryptBase64(unencrypted []byte, key []byte, iv []byte) (encrypted string, err error) {
	d, err := AESEncrypt(unencrypted, key, iv)
	if err != nil {
		return
	}
	return base64.StdEncoding.EncodeToString(d), nil
}

// AESNonceEncryptBase64 aes encrypt data with given key and random bytes as IV.
//Data will be padding with PKCS7Padding
//Random IV will prefix encryped data
//return base64 encoded encrypted data and any error if raisd.
func AESNonceEncryptBase64(unencrypted []byte, key []byte) (encrypted string, err error) {
	d, err := AESNonceEncrypt(unencrypted, key)
	if err != nil {
		return
	}
	return base64.StdEncoding.EncodeToString(d), nil
}

//AESDecrypt decrypt data with given key and iv.
//Data will be unpadding with PKCS7Unpadding.
//Return decrypted data and any error if raised.
func AESDecrypt(encrypted []byte, key []byte, iv []byte) (decrypted []byte, err error) {
	defer func() {
		r := recover()
		if r != nil {
			err = r.(error)
		}
	}()
	cryptKey := formatKey(key, aes.BlockSize)
	block, err := aes.NewCipher(cryptKey)
	if err != nil {
		return
	}
	crypter := cipher.NewCBCDecrypter(block, iv)
	data := make([]byte, len(encrypted))
	crypter.CryptBlocks(data, encrypted)
	decrypted = PKCS7Unpadding(data)
	return
}

//AESNonceDecrypt decrypt data with given key.
//IV will load form first bytes of data.
//Data will be unpadding with PKCS7Unpadding.
//Return decrypted data and any error if raised.
func AESNonceDecrypt(encrypted []byte, key []byte) (decrypted []byte, err error) {
	defer func() {
		r := recover()
		if r != nil {
			err = r.(error)
		}
	}()
	return AESDecrypt(encrypted[IVSize:], key, encrypted[:IVSize])
}

//AESDecryptBase64 decrypt base64 encoded data with given key and iv.
//Data will be unpadding with PKCS7Unpadding.
//Return decrypted data and any error if raised.
func AESDecryptBase64(encrypted string, key []byte, iv []byte) (decrypted []byte, err error) {
	d, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return
	}
	return AESDecrypt(d, key, iv)
}

//AESNonceDecryptBase64 decrypt base64 encoded data with given key.
//IV will load form first bytes of data.
//Data will be unpadding with PKCS7Unpadding.
//Return decrypted data and any error if raised.
func AESNonceDecryptBase64(encrypted string, key []byte) (decrypted []byte, err error) {
	d, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return
	}
	return AESNonceDecrypt(d, key)
}

//AESTokenMarshaler token marshaler which crypt data with AES
//Return error if raised
func AESTokenMarshaler(s *ClientDriver, ts *Session) (err error) {
	var data []byte

	ts.Nonce = make([]byte, clientStoreNonceSize)
	_, err = rand.Read(ts.Nonce)
	if err != nil {
		return
	}
	data, err = ts.Marshal()
	if err != nil {
		return
	}
	ts.token, err = AESNonceEncryptBase64(data, s.Key)
	return
}

//AESTokenUnmarshaler token unmarshaler which crypt data with AES
//Return error if raised
func AESTokenUnmarshaler(s *ClientDriver, v *Session) (err error) {
	var data []byte
	data, err = AESNonceDecryptBase64(v.token, s.Key)
	if err != nil {
		return ErrDataNotFound
	}
	err = v.Unmarshal(v.token, data)
	if err != nil {
		return ErrDataNotFound
	}
	return nil
}

//ClientDriver ClientDriver is the stuct store token data in Client side.
type ClientDriver struct {
	Key              []byte                              //Crypt key
	TokenMarshaler   func(*ClientDriver, *Session) error //Marshler data to Session.token
	TokenUnmarshaler func(*ClientDriver, *Session) error //Unmarshler data from Session.token
}

//NewClientDriver New create a new client side token store with given key and token lifetime.
//Key the key used to encrpty data
//TokenLifeTime is the token initial expired tome.
//Return a new token store.
//All other property of the store can be set after creation.
func NewClientDriver() *ClientDriver {
	return &ClientDriver{
		TokenMarshaler:   AESTokenMarshaler,
		TokenUnmarshaler: AESTokenUnmarshaler,
	}
}

//MustClientStore create new client store with given  key and ttl.
//Return store created.
//Panic if any error raised.
func MustClientStore(key []byte, TokenLifetime time.Duration) *Store {
	driver := NewClientDriver()
	oc := NewClientDriverOptionConfig()
	oc.Key = key
	err := driver.Init(oc)
	if err != nil {
		panic(err)
	}
	store := New()
	store.Marshaler, err = cache.NewMarshaler(DefaultMarshaler)
	if err != nil {
		panic(err)
	}
	soc := NewOptionConfig()
	soc.Driver = driver
	soc.TokenLifetime = TokenLifetime
	err = store.Init(soc)
	if err != nil {
		panic(err)
	}
	return store
}

//GetSessionToken Get the token string from token data.
//Return token and any error raised.
func (s *ClientDriver) GetSessionToken(ts *Session) (token string, err error) {
	err = ts.Save()
	return ts.token, err
}

//Init init client driver with given option.
//Return any error if raised.
func (s *ClientDriver) Init(option ClientDriverOption) error {
	return option.ApplyTo(s)
}

//GenerateToken generate new token name with given prefix.
//Return the new token name and error.
func (s *ClientDriver) GenerateToken(prefix string) (token string, err error) {
	return clientStoreNewToken, nil

}

//Load Load Session form the Session.token.
//Return any error if raised
func (s *ClientDriver) Load(v *Session) (err error) {
	err = s.TokenUnmarshaler(s, v)
	if err != nil {
		return err
	}
	return
}

//Save Save Session if necessary.
//Return any error raised.
func (s *ClientDriver) Save(ts *Session, ttl time.Duration) (err error) {
	ts.oldToken = ts.token
	err = s.TokenMarshaler(s, ts)
	if err != nil {
		return
	}
	if ts.oldToken != ts.token {
		ts.tokenChanged = true
	}
	return
}

//Delete delete the token with given name.
//Return any error if raised.
func (s *ClientDriver) Delete(token string) (bool, error) {
	return false, nil
}

//Close Close cachestore and return any error if raised
func (s *ClientDriver) Close() error {
	return nil
}
func (s *ClientDriver) DynamicToken() bool {
	return true
}

// PKCS7Padding padding data as  PKCS7
// Reference http://blog.studygolang.com/167.html
func PKCS7Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	d := make([]byte, padding+len(data))
	copy(d, data)
	copy(d[len(data):], padtext)
	return d

}

// PKCS7Unpadding unpadding data as  PKCS7
// Reference http://blog.studygolang.com/167.html
func PKCS7Unpadding(data []byte) []byte {
	length := len(data)
	unpadding := int(data[length-1])
	d := make([]byte, length-unpadding)
	copy(d, data)
	return d
}
