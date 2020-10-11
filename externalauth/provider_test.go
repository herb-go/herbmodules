package auth

import "net/http"

func newTestDriver(data Profile) *testDriver {
	return &testDriver{
		data: data,
	}
}

type testCreator struct {
	Profile
}

func (c *testCreator) Create() Driver {
	return newTestDriver(c.Profile)
}

type testDriver struct {
	data Profile
}

func (d *testDriver) ExternalLogin(provider *Provider, w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, provider.AuthURL(), 301)
}
func (d *testDriver) AuthRequest(provider *Provider, r *http.Request) (*Result, error) {
	result := provider.Auth.MustGetResult(r)
	result.Account = d.data[ProfileIndexID][0]
	result.Data = d.data
	return result, nil
}

func newTestFailDriver(data Profile) *testFailDriver {
	return &testFailDriver{
		data: data,
	}
}

type testFailDriver struct {
	data Profile
}

func (d *testFailDriver) ExternalLogin(provider *Provider, w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, provider.AuthURL(), 301)
}
func (d *testFailDriver) AuthRequest(provider *Provider, r *http.Request) (*Result, error) {
	result := provider.Auth.MustGetResult(r)
	return result, ErrAuthParamsError
}
