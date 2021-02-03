package payloadauthorizer

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

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
func (a *Authorizer) With(l protecter.PolicyLoader) *Authorizer {
	a.PolicyLoader = l
	return a
}
func (a *Authorizer) WithAny(pls ...protecter.PolicyLoader) *Authorizer {
	a.PolicyLoader = Any(pls...)
	return a
}
func (a *Authorizer) WithAll(pls ...protecter.PolicyLoader) *Authorizer {
	a.PolicyLoader = All(pls...)
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

func Any(pls ...protecter.PolicyLoader) protecter.PolicyLoader {
	return protecter.PolicyLoaderFunc(func(r *http.Request) (role.Policy, error) {
		var err error
		ps := make([]role.Policy, len(pls))
		for k := range pls {
			ps[k], err = pls[k].LoadPolicy(r)
			if err != nil {
				return nil, err
			}
		}
		return role.Any(ps...), nil
	})
}

func Not(pl protecter.PolicyLoader) protecter.PolicyLoader {
	return protecter.PolicyLoaderFunc(func(r *http.Request) (role.Policy, error) {
		p, err := pl.LoadPolicy(r)
		if err != nil {
			return nil, err
		}
		return role.Not(p), nil
	})
}

func All(pls ...protecter.PolicyLoader) protecter.PolicyLoader {
	return protecter.PolicyLoaderFunc(func(r *http.Request) (role.Policy, error) {
		var err error
		ps := make([]role.Policy, len(pls))
		for k := range pls {
			ps[k], err = pls[k].LoadPolicy(r)
			if err != nil {
				return nil, err
			}
		}
		return role.All(ps...), nil
	})
}

func MustParse(str string) protecter.PolicyLoader {
	p, err := roleparser.Parse(str)
	if err != nil {
		panic(err)
	}
	return protecter.PolicyLoaderFunc(func(r *http.Request) (role.Policy, error) {
		return p, nil
	})
}

func makeToken(index int) string {
	return fmt.Sprintf("{{%d}}", index)
}

type replacer struct {
	token string
	load  func(r *http.Request) (string, error)
}

func MustParseWith(pattern string, paramsloaders ...func(r *http.Request) (string, error)) protecter.PolicyLoader {
	var replacers = []*replacer{}
	for k := range paramsloaders {
		replacers = append(replacers, &replacer{
			token: makeToken(k),
			load:  paramsloaders[k],
		})
	}
	return protecter.PolicyLoaderFunc(func(r *http.Request) (role.Policy, error) {
		var replacements = make([]string, 0, len(paramsloaders)*2)
		for k := range replacers {
			value, err := replacers[k].load(r)
			if err != nil {
				return nil, err
			}
			replacements = append(replacements, replacers[k].token, url.PathEscape(value))
		}
		return roleparser.Parse(strings.NewReplacer(replacements...).Replace(pattern))
	})
}
