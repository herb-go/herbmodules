package healthcheck

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

var testcode = Code(1)
var testhealthychecker = func() (Status, *Info) {
	return StatusHealthy, nil
}
var testhealthydatachecker = func() (Status, *Info) {
	return StatusHealthy, NewInfo().WithMsg("testmsg").WithCode(&testcode, "test")
}
var testwarningchecker = func() (Status, *Info) {
	return StatusWarning, nil
}
var testwarningdatachecker = func() (Status, *Info) {
	return StatusWarning, NewInfo().WithMsg("testmsg").WithCode(&testcode, "test")
}

var testerrorchecker = func() (Status, *Info) {
	return StatusError, nil
}
var testerrordatachecker = func() (Status, *Info) {
	return StatusError, NewInfo().WithMsg("testmsg").WithCode(&testcode, "test")
}

func TestCheck(t *testing.T) {
	Reset()
	defer Reset()
	s := httptest.NewServer(Hanlder)
	defer s.Close()
	resp, err := http.DefaultClient.Get(s.URL)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != StatusCodeHealthy {
		t.Fatal(resp)
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	result := NewResult()
	err = json.Unmarshal(data, result)
	if err != nil {
		panic(err)
	}
	if result.Status != StatusHealthy || result.Msgs != nil || result.Warnings != nil || result.Errors != nil {
		t.Fatal(result)
	}
	resp.Body.Close()

	OnCheck(testhealthychecker)

	resp, err = http.DefaultClient.Get(s.URL)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != StatusCodeHealthy {
		t.Fatal(resp)
	}
	data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	result = NewResult()
	err = json.Unmarshal(data, result)
	if err != nil {
		panic(err)
	}
	if result.Status != StatusHealthy || result.Msgs != nil || result.Warnings != nil || result.Errors != nil {
		t.Fatal(result)
	}
	resp.Body.Close()

	OnCheck(testhealthydatachecker)
	resp, err = http.DefaultClient.Get(s.URL)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != StatusCodeHealthy {
		t.Fatal(resp)
	}
	data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	result = NewResult()
	err = json.Unmarshal(data, result)
	if err != nil {
		panic(err)
	}
	if result.Status != StatusHealthy || len(*result.Msgs) != 1 || result.Warnings != nil || result.Errors != nil {
		t.Fatal(result)
	}
	m := (*result.Msgs)[0]
	if m.Msg != "testmsg" || *m.Code != testcode || m.Data == nil {
		t.Fatal(m)
	}

	resp.Body.Close()

	OnCheck(testwarningchecker)

	resp, err = http.DefaultClient.Get(s.URL)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != StatusCodeWarning {
		t.Fatal(resp)
	}
	data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	result = NewResult()
	err = json.Unmarshal(data, result)
	if err != nil {
		panic(err)
	}
	if result.Status != StatusWarning || len(*result.Msgs) != 1 || result.Warnings != nil || result.Errors != nil {
		t.Fatal(result)
	}
	resp.Body.Close()

	OnCheck(testwarningdatachecker)
	resp, err = http.DefaultClient.Get(s.URL)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != StatusCodeWarning {
		t.Fatal(resp)
	}
	data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	result = NewResult()
	err = json.Unmarshal(data, result)
	if err != nil {
		panic(err)
	}
	if result.Status != StatusWarning || len(*result.Msgs) != 1 || len(*result.Warnings) != 1 || result.Errors != nil {
		t.Fatal(result)
	}
	m = (*result.Warnings)[0]
	if m.Msg != "testmsg" || *m.Code != testcode || m.Data == nil {
		t.Fatal(m)
	}

	resp.Body.Close()

	OnCheck(testerrorchecker)

	resp, err = http.DefaultClient.Get(s.URL)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != StatusCodeError {
		t.Fatal(resp)
	}
	data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	result = NewResult()
	err = json.Unmarshal(data, result)
	if err != nil {
		panic(err)
	}
	if result.Status != StatusError || len(*result.Msgs) != 1 || len(*result.Warnings) != 1 || result.Errors != nil {
		t.Fatal(result)
	}
	resp.Body.Close()

	OnCheck(testerrordatachecker)
	resp, err = http.DefaultClient.Get(s.URL)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != StatusCodeError {
		t.Fatal(resp)
	}
	data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	result = NewResult()
	err = json.Unmarshal(data, result)
	if err != nil {
		panic(err)
	}
	if result.Status != StatusError || len(*result.Msgs) != 1 || len(*result.Warnings) != 1 || len(*result.Errors) != 1 {
		t.Fatal(result)
	}
	m = (*result.Errors)[0]
	if m.Msg != "testmsg" || *m.Code != testcode || m.Data == nil {
		t.Fatal(m)
	}

	resp.Body.Close()
}
