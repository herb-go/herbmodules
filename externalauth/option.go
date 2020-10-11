package auth

import "net/url"

//Option auth option interface
type Option interface {
	ApplyTo(*Auth) error
}

//OptionFunc auth option function interface.
type OptionFunc func(*Auth) error

//ApplyTo apply option function to auth.
func (i OptionFunc) ApplyTo(a *Auth) error {
	return i(a)
}

//OptionCommon auth service option with given hostpath and session.
//Params hostpath shold be form "https://www.example.com/auth"
func OptionCommon(hostpath string, Session Session) OptionFunc {
	return func(a *Auth) error {
		u, err := url.Parse(hostpath)
		if err != nil {
			return err
		}
		*a = Auth{
			LoginPrefix:    DefaultLoginPrefix,
			AuthPrefix:     DefaultAuthPrefix,
			Session:        Session,
			Host:           u.Scheme + "://" + u.Host,
			Path:           u.Path,
			NotFoundAction: DefaultNotFoundAction,
		}
		return nil
	}
}
