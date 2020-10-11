package auth

import "net/http"

//Driver auth provider driver ineterface.
type Driver interface {
	//ExternalLogin action which redirect to login page.
	ExternalLogin(provider *Provider, w http.ResponseWriter, r *http.Request)
	//AuthRequest auth provider response request.
	//Return auth result and any error raised.
	//If auth request params is not correct,error ErrAuthParamsError will be returned.
	AuthRequest(provider *Provider, r *http.Request) (*Result, error)
}

// DriverCreator driver  creator interface.
type DriverCreator interface {
	Create() Driver
}

//Provider auth provider mian struct
type Provider struct {
	//Driver auth driver
	Driver Driver
	//Auth auth service which provider belongs to.
	Auth *Auth
	//Keyword provider keyword.
	Keyword string
}

//Login login action.
func (p *Provider) Login(w http.ResponseWriter, r *http.Request) {
	p.Driver.ExternalLogin(p, w, r)
}

//AuthRequest auth provider response request.
//Return auth result and any error raised.
//If auth request params is not correct,error ErrAuthParamsError will be returned.
func (p *Provider) AuthRequest(r *http.Request) (*Result, error) {
	return p.Driver.AuthRequest(p, r)
}

//LoginURL return provider's login url.
func (p *Provider) LoginURL() string {
	return p.Auth.Host + p.Auth.Path + p.Auth.LoginPrefix + p.Keyword
}

//AuthURL return provider's auth url.
func (p *Provider) AuthURL() string {
	return p.Auth.Host + p.Auth.Path + p.Auth.AuthPrefix + p.Keyword
}
