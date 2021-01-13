package payloadauthorizer

import (
	"net/http"

	"github.com/herb-go/herbsecurity/authorize/role"
	"github.com/herb-go/herbsecurity/authorize/role/roleparser"

	"github.com/herb-go/herbsecurity/authority"

	"github.com/herb-go/herbmodules/protecter"
)

var DefaultPolicyLoader = protecter.PolicyLoaderFunc(func(r *http.Request) (role.Policy, error) {
	return role.Deny, nil
})
var DefaultOnDeny = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(403), 403)
})

type Authorizer struct {
	Key          protecter.Key
	PolicyLoader protecter.PolicyLoader
	OnDeny       http.Handler
}

func (a *Authorizer) WithOnDeny(d http.Handler) *Authorizer {
	a.OnDeny = d
	return a
}
func (a *Authorizer) WithKey(key protecter.Key) *Authorizer {
	a.Key = key
	return a
}
func (a *Authorizer) WithPolicy(r role.Policy) *Authorizer {
	a.PolicyLoader = protecter.RolePolicyLoader(r)
	return a
}
func (a *Authorizer) WithPolicyLoader(l protecter.PolicyLoader) *Authorizer {
	a.PolicyLoader = l
	return a
}
func (a *Authorizer) ServeMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	auth := a.Key.LoadAuth(r)
	if auth != nil {
		rolesstr := auth.Payloads().LoadString(authority.PayloadRoles)
		role, err := roleparser.Parse(rolesstr)
		if err != nil {
			panic(err)
		}
		policy, err := a.PolicyLoader.LoadPolicy(r)
		if err != nil {
			panic(err)
		}
		ok, err := policy.Authorize(role)
		if err != nil {
			panic(err)
		}
		if ok {
			next(w, r)
			return
		}
	}
	a.OnDeny.ServeHTTP(w, r)
}

func New() *Authorizer {
	return &Authorizer{
		Key:          protecter.DefaultKey,
		PolicyLoader: DefaultPolicyLoader,
		OnDeny:       DefaultOnDeny,
	}
}
