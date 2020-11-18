package httpsession

import (
	"bytes"
	"testing"

	"github.com/herb-go/herbdata-drivers/kvdb-drivers/freecachedb"
	"github.com/herb-go/herbdata/kvdb"

	"github.com/herb-go/herbdata"
)

func TestKVEngine(t *testing.T) {
	var err error
	var data []byte
	var e *KVEngine
	var token string
	var newtoken string
	var db = kvdb.New()
	driver, err := (&freecachedb.Config{Size: 500000}).CreateDriver()
	if err != nil {
		panic(err)
	}
	db.Driver = driver
	e = &KVEngine{
		TokenSize: 256,
		Timeout:   0,
		Database:  db,
	}
	if e.DynamicToken() {
		t.Fatal()
	}
	if e.SessionTimeout() != 0 {
		t.Fatal()
	}
	token, err = e.NewToken()
	if token == TokenEmpty || err != nil {
		t.Fatal(token, err)
	}
	data = []byte("data")
	token, err = e.SaveToken(token, data, 1200)
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
	newtoken, data, err = e.LoadToken(token)
	if err != herbdata.ErrNotFound {
		t.Fatal(newtoken, data, err)
	}
}

func TestKVEngineNotFound(t *testing.T) {
	var err error
	var e *KVEngine
	var db = kvdb.New()
	driver, err := (&freecachedb.Config{Size: 500000}).CreateDriver()
	if err != nil {
		panic(err)
	}
	db.Driver = driver
	e = &KVEngine{
		TokenSize: 256,
		Timeout:   120,
		Database:  db,
	}
	_, _, err = e.LoadToken("!notexist")
	if err != herbdata.ErrNotFound {
		t.Fatal()
	}
	_, _, err = e.LoadToken("")
	if err != herbdata.ErrNotFound {
		t.Fatal()
	}
}
