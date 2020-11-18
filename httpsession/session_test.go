package httpsession

import (
	"bytes"
	"testing"

	"github.com/herb-go/herbdata"
)

func TestSession(t *testing.T) {
	s := newSession()
	s.MarkAsStarted()
	data, err := s.Get([]byte("testkey"))
	if err != herbdata.ErrNotFound {
		t.Fatal(data, err)
	}
	err = s.Set([]byte("testkey"), []byte("testdata"))
	if err != nil {
		t.Fatal(err)
	}
	data, err = s.Get([]byte("testkey"))
	if err != nil || bytes.Compare(data, []byte("testdata")) != 0 {
		t.Fatal(data, err)
	}
	err = s.Delete([]byte("testkey"))
	if err != nil {
		t.Fatal(err)
	}
	data, err = s.Get([]byte("testkey"))
	if err != herbdata.ErrNotFound {
		t.Fatal(data, err)
	}
	s = newSession()
	ok := s.Updated()
	if ok != false {
		t.Fatal(ok)
	}
	s.MarkAsUpdated()
	ok = s.Updated()
	if ok != true {
		t.Fatal(ok)
	}
	s = newSession()
	ok = s.Temporay()
	if ok != false {
		t.Fatal(ok)
	}
	if s.Updated() {
		t.Fatal()
	}
	s.MarkAsTemporay()
	ok = s.Temporay()
	if ok != true {
		t.Fatal(ok)
	}
	if !s.Updated() {
		t.Fatal()
	}
	s = newSession()
	s.MarkAsStarted()
	if s.Updated() {
		t.Fatal()
	}
	err = s.Set([]byte("testkey"), []byte("testdata"))
	if err != nil {
		t.Fatal(err)
	}
	if !s.Updated() {
		t.Fatal()
	}
	s = newSession()
	s.MarkAsStarted()
	if s.Updated() {
		t.Fatal()
	}
	err = s.Delete([]byte("testkey"))
	if err != nil {
		t.Fatal(err)
	}
	if !s.Updated() {
		t.Fatal()
	}
	sd := newSessionData()
	sd.CreatedAt = 1
	sd.ExpiredAt = 2
	sd.Data = append(sd.Data, &SessionValue{Key: "test", Value: []byte("testdata")})
	s = newSession()
	s.MarkAsStarted()
	data, err = s.Load("test")
	if err != herbdata.ErrNotFound {
		t.Fatal(data, err)
	}
	if *s.expiredAt != 0 || *s.createdAt != 0 {
		t.Fatal(s)
	}
	s.setdata(sd)
	data, err = s.Load("test")
	if err != nil || bytes.Compare(data, []byte("testdata")) != 0 {
		t.Fatal(data, err)
	}
	if *s.expiredAt != 2 || *s.createdAt != 1 {
		t.Fatal(s)
	}
	sd2 := s.getdata()
	if sd2.ExpiredAt != sd.ExpiredAt || sd2.CreatedAt != sd.CreatedAt ||
		len(sd2.Data) != len(sd.Data) || sd2.Data[0].Key != sd.Data[0].Key ||
		bytes.Compare(sd2.Data[0].Value, sd.Data[0].Value) != 0 {
		t.Fatal(sd, sd2)
	}
	s.MarkAsTemporay()
	sd2 = s.getdata()
	if !sd2.Temporay {
		t.Fatal()
	}
	s = newSession()
	if s.Temporay() {
		t.Fatal()
	}
	s.setdata(sd2)
	if !s.Temporay() {
		t.Fatal()
	}
}

func TestNotStarted(t *testing.T) {
	var err error
	s := newSession()
	if s.Started() {
		t.Fatal()
	}
	err = s.Set([]byte("test"), []byte("data"))
	if err != ErrSessionNotStarted {
		t.Fatal()
	}
	_, err = s.Get([]byte("test"))
	if err != ErrSessionNotStarted {
		t.Fatal()
	}
	err = s.Delete([]byte("test"))
	if err != ErrSessionNotStarted {
		t.Fatal()
	}
	s.MarkAsStarted()
	if !s.Started() {
		t.Fatal()
	}
	err = s.Set([]byte("test"), []byte("data"))
	if err != nil {
		t.Fatal()
	}
	_, err = s.Get([]byte("test"))
	if err != nil {
		t.Fatal()
	}
	err = s.Delete([]byte("test"))
	if err != nil {
		t.Fatal()
	}
}
