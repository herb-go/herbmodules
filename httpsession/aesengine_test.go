package httpsession

import (
	"bytes"
	"testing"

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
	e = newTestAESEngine()
	_, _, err = e.LoadToken("12345")
	if err != herbdata.ErrNotFound {
		t.Fatal(err)
	}
}

func TestAESEngine(t *testing.T) {
	var err error
	var data []byte
	var e *AESEngine
	var token string
	var newtoken string
	e = newTestAESEngine()
	if !e.DynamicToken() {
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
	err = e.RevokeToken(token)
	if err != nil {
		t.Fatal()
	}
}

func TestAESEngineNotFound(t *testing.T) {
	var err error
	var e *AESEngine
	e = newTestAESEngine()
	_, _, err = e.LoadToken("!notexist")
	if err != herbdata.ErrNotFound {
		t.Fatal()
	}
	_, _, err = e.LoadToken("")
	if err != herbdata.ErrNotFound {
		t.Fatal()
	}
}
