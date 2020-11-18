package httpsession

import (
	"bytes"
	"testing"
	"time"

	"github.com/herb-go/herbdata"
)

var _ Engine = &AESEngine{}

func newTestAESEngine() *AESEngine {
	return &AESEngine{
		Secret: []byte("secret"),
	}
}
func TestAESEngineParseFail(t *testing.T) {
	var err error
	var e *AESEngine
	var token string
	e = newTestAESEngine()
	token, err = AESNonceEncryptBase64([]byte("12345"), e.Secret)
	if err != nil {
		t.Fatal()
	}
	_, _, err = e.LoadToken(token)
	if err != herbdata.ErrNotFound {
		t.Fatal()
	}
}
func TestAESEngine(t *testing.T) {
	var err error
	var data []byte
	var e *AESEngine
	var token string
	var newtoken string
	var lastactive int64
	e = newTestAESEngine()
	e.Timeout = 2
	if !e.DynamicToken() {
		t.Fatal()
	}
	if e.SessionTimeout() != 2 {
		t.Fatal()
	}
	token, err = e.NewToken()
	if token == TokenEmpty || err != nil {
		t.Fatal(token)
	}
	data = []byte("data")
	token, err = e.SaveToken(TokenEmpty, data, 0)
	if token == TokenEmpty || err != nil {
		t.Fatal(token)
	}
	newtoken, data, err = e.LoadToken(token)
	if newtoken == token || bytes.Compare(data, []byte("data")) != 0 || err != nil {
		t.Fatal(newtoken, data, err)
	}
	lastactive, err = e.TokenLastActive(newtoken)
	if time.Now().Unix()-lastactive > 1 || err != nil {
		t.Fatal()
	}
	token = newtoken
	time.Sleep(3 * time.Second)
	newtoken, data, err = e.LoadToken(token)
	if err != herbdata.ErrNotFound {
		t.Fatal(newtoken, data, err)
	}
	lastactive, err = e.TokenLastActive(token)
	if err != herbdata.ErrNotFound {
		t.Fatal(newtoken, data, err)
	}

}

func TestAESEngineWithZeroTimeout(t *testing.T) {
	var err error
	var data []byte
	var e *AESEngine
	var token string
	var newtoken string
	var lastactive int64
	e = newTestAESEngine()
	e.Timeout = 0
	if !e.DynamicToken() {
		t.Fatal()
	}
	if e.SessionTimeout() != 0 {
		t.Fatal()
	}
	token, err = e.NewToken()
	if token != TokenEmpty || err != nil {
		t.Fatal(token)
	}
	data = []byte("data")
	token, err = e.SaveToken(TokenEmpty, data, 0)
	if token == TokenEmpty || err != nil {
		t.Fatal(token, err)
	}
	newtoken, data, err = e.LoadToken(token)
	if newtoken != token || bytes.Compare(data, []byte("data")) != 0 || err != nil {
		t.Fatal(newtoken, data, err)
	}
	lastactive, err = e.TokenLastActive(newtoken)
	if lastactive != 0 || err != nil {
		t.Fatal()
	}
	err = e.RevokeToken(token)
	if err != nil {
		t.Fatal()
	}
}

func TestAESEngineNotFound(t *testing.T) {
	var err error
	var e *AESEngine
	e = newTestAESEngine()
	e.Timeout = 10
	_, _, err = e.LoadToken("!notexist")
	if err != herbdata.ErrNotFound {
		t.Fatal()
	}
	_, _, err = e.LoadToken("")
	if err != herbdata.ErrNotFound {
		t.Fatal()
	}
	_, err = e.TokenLastActive("!notexist")
	if err != herbdata.ErrNotFound {
		t.Fatal()
	}
	_, err = e.TokenLastActive("")
	if err != herbdata.ErrNotFound {
		t.Fatal()
	}
}
