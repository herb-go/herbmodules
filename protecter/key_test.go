package protecter

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/herb-go/herb/middleware"
	"github.com/herb-go/herbsecurity/authorize/role"

	"github.com/herb-go/herbsecurity/authority"

	"github.com/herb-go/herbsecurity/authority/credential"
)

var okHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ok"))
})
var testHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	id, err := DefaultKey.IdentifyRequest(r)
	if err != nil {
		panic(err)
	}
	w.Write([]byte(id))
})
var testIDHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(LoadAuth(r).Principal()))
})
var credentialerAppID = CredentialerFunc(func(r *http.Request) credential.CredentialSource {
	return credential.New().WithName(credential.NameAuthority).WithValue([]byte(r.Header.Get("appid")))
})

var credentialerToken = CredentialerFunc(func(r *http.Request) credential.CredentialSource {
	return credential.New().WithName(credential.NamePassphrase).WithValue([]byte(r.Header.Get("token")))
})

var notfound = http.NotFoundHandler()

var testProtecter = New().
	WithOnFail(notfound).
	WithCredentialers(credentialerAppID, credentialerToken).
	WithAuthenticator(
		credential.AuthenticatorFunc(func(c credential.Credentials) (*authority.Auth, error) {
			auth, err := c.Get(credential.NameAuthority)
			if err != nil {
				return nil, err
			}
			pass, err := c.Get(credential.NamePassphrase)
			if err != nil {
				return nil, err
			}
			if string(auth) == "testappid" && string(pass) == "testtoken" {
				return authority.NewAuth("testappid"), nil
			}
			return nil, nil
		},
			credential.NameAuthority,
			credential.NamePassphrase,
		),
	)

func TestForbidden(t *testing.T) {
	s := httptest.NewServer(ProtectWith(ForbiddenProtecter, testHandler))
	defer s.Close()
	req, err := http.NewRequest("GET", s.URL, nil)
	if err != nil {
		panic(err)
	}
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		panic(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 403 {
		t.Fatal(resp)
	}
}

func TestSuccess(t *testing.T) {
	s := httptest.NewServer(ProtectWith(NotWorkingProtecter, testHandler))
	defer s.Close()
	req, err := http.NewRequest("GET", s.URL, nil)
	if err != nil {
		panic(err)
	}
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		panic(err)
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Fatal(resp)
	}
	if string(data) != "notworking" {
		t.Fatal()
	}
}

func TestNil(t *testing.T) {
	s := httptest.NewServer(ProtectWith(nil, testHandler))
	defer s.Close()
	req, err := http.NewRequest("GET", s.URL, nil)
	if err != nil {
		panic(err)
	}
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		panic(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 403 {
		t.Fatal(resp)
	}
}

func TestAuthFail(t *testing.T) {
	s := httptest.NewServer(ProtectWith(testProtecter, testHandler))
	defer s.Close()
	req, err := http.NewRequest("GET", s.URL, nil)
	if err != nil {
		panic(err)
	}
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		panic(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 404 {
		t.Fatal(resp)
	}

}

func TestAuthSuccess(t *testing.T) {
	s := httptest.NewServer(ProtectWith(testProtecter, testIDHandler))
	defer s.Close()
	req, err := http.NewRequest("GET", s.URL, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Add("appid", "testappid")
	req.Header.Add("token", "testtoken")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Fatal(resp)
	}
	if string(data) != "testappid" {
		t.Fatal(string(data))
	}
}

func TestMiddlewareFail(t *testing.T) {
	m := ProtectMiddleware(nil)
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m(w, r, testHandler)
	}))
	defer s.Close()
	resp, err := http.Get(s.URL)
	if err != nil {
		panic(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 403 {
		t.Fatal(resp)
	}
}

func TestMiddlewareSuccess(t *testing.T) {
	m := ProtectMiddleware(NotWorkingProtecter)
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m(w, r, testHandler)
	}))
	defer s.Close()
	resp, err := http.Get(s.URL)
	if err != nil {
		panic(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Fatal(resp)
	}
}

func TestRoles(t *testing.T) {
	loader1 := RoleRolesLoader(role.New("role1"))
	loader2 := RoleRolesLoader(role.New("role2"))
	policy1 := RolePolicyLoader(role.New("role1"))
	policy2 := RolePolicyLoader(role.New("role2"))
	mux := &http.ServeMux{}
	mux.Handle("/l1p1", middleware.New().Use(
		RolesMiddleware(loader1),
		PolicyMiddleware(policy1),
	).Handle(okHandler),
	)
	mux.Handle("/l1l2p1", middleware.New().Use(
		RolesMiddleware(loader1, loader2),
		PolicyMiddleware(policy1),
	).Handle(okHandler),
	)
	mux.Handle("/l1p1p2", middleware.New().Use(
		RolesMiddleware(loader1),
		PolicyMiddleware(policy1, policy2),
	).Handle(okHandler),
	)
	mux.Handle("/l1l2p1p2", middleware.New().Use(
		RolesMiddleware(loader1, loader2),
		PolicyMiddleware(policy1, policy2),
	).Handle(okHandler),
	)
	s := httptest.NewServer(mux)
	defer s.Close()
	resp, err := http.Get(s.URL + "/l1p1")
	if err != nil {
		panic(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Fatal(resp.StatusCode)
	}
	resp, err = http.Get(s.URL + "/l1l2p1")
	if err != nil {
		panic(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Fatal(resp.StatusCode)
	}
	resp, err = http.Get(s.URL + "/l1p1p2")
	if err != nil {
		panic(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 403 {
		t.Fatal(resp.StatusCode)
	}
	resp, err = http.Get(s.URL + "/l1l2p1p2")
	if err != nil {
		panic(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Fatal(resp.StatusCode)
	}
}
