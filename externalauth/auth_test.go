package auth

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/herb-go/herbmodules/externalauth/testsession"
)

func TestMisc(t *testing.T) {
	Auth := New()
	bs, err := Auth.RandToken(32)
	if err != nil {
		t.Fatal(err)
	}
	if len(bs) != 32 {
		t.Error(t)
	}
}
func TestServer(t *testing.T) {
	Auth := New()

	var successAction = func(w http.ResponseWriter, r *http.Request) {
		result := Auth.MustGetResult(r)
		data, err := json.Marshal(result)
		if err != nil {
			panic(err)
		}
		w.Write(data)
	}
	var failAction = func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, http.StatusText(400), 400)
	}
	var Session = &testsession.TestCookieSession{
		Name: "herbgoexternalauthlogin",
		Path: "/",
	}
	mux := http.StripPrefix("/auth", http.HandlerFunc(Auth.Serve(successAction)))
	server := httptest.NewServer(mux)
	err := Auth.Init(OptionCommon(server.URL+"/auth", Session))
	data := Profile{
		ProfileIndexID: []string{"test"},
	}
	Auth.MustRegisterWithCreator("test", &testCreator{data})
	provider := Auth.MustGetProvider("test")
	resp, err := http.Get(provider.LoginURL())
	if err != nil {
		t.Fatal(err)
	}

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Error(resp.StatusCode)
	}
	result := &Result{}
	err = json.Unmarshal(content, result)
	if err != nil {
		t.Fatal(err)
	}
	if result.Account != data[ProfileIndexID][0] {
		t.Error(result.Account)
	}
	if result.Data[ProfileIndexID][0] != data[ProfileIndexID][0] {
		t.Error(result.Data[ProfileIndexID][0])
	}

	resp, err = http.Get(server.URL + "/auth/login/notexist")
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 404 {
		t.Error(resp.StatusCode)
	}

	_, err = Auth.RegisterProvider("fail", newTestFailDriver(data))
	if err != nil {
		t.Fatal(err)
	}
	provider, err = Auth.GetProvider("fail")
	if err != nil {
		t.Fatal(err)
	}
	resp, err = http.Get(provider.LoginURL())
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 404 {
		t.Error(resp.StatusCode)
	}
	Auth.NotFoundAction = failAction
	resp, err = http.Get(provider.LoginURL())
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 400 {
		t.Error(resp.StatusCode)
	}
	pm := newTestProviderManager()
	Auth.SetProviderManager(pm)
	Auth.MustRegisterProvider("test", newTestDriver(data))
	provider = Auth.MustGetProvider("test")
	resp, err = http.Get(provider.LoginURL())
	if err != nil {
		t.Fatal(err)
	}

	content, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Error(resp.StatusCode)
	}
	result = &Result{}
	err = json.Unmarshal(content, result)
	if err != nil {
		t.Fatal(err)
	}
	if result.Account != data[ProfileIndexID][0] {
		t.Error(result.Account)
	}
	if result.Data[ProfileIndexID][0] != data[ProfileIndexID][0] {
		t.Error(result.Data[ProfileIndexID][0])
	}

}
