package herbsession

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
)

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
