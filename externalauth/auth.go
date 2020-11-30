package auth

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

//DefaultLoginPrefix default login action prefix
const DefaultLoginPrefix = "/login/"

//DefaultAuthPrefix default auth action prefix
const DefaultAuthPrefix = "/auth/"

//ErrAuthParamsError error raised when auth prarms error.
//Should be rasied in auth driver's AuthRequest method.
var ErrAuthParamsError = errors.New("external auth params error")

//TokenMask used to generate random bytes.
var TokenMask = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789_-")

//ContextName auth result context name type
type ContextName string

//ResultContextName context name for result in request.
const ResultContextName = ContextName("authresult")

//DefaultNotFoundAction default not found action execute when provider not found.
func DefaultNotFoundAction(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
}

//Auth auth service main struct
type Auth struct {
	Debug           bool
	providerManager ProviderManager
	//Host auth service host.For example,https://www.example.com
	Host string
	//Path auth servcei path.For example,"/auth" in https://www.example.com/auth
	Path string
	//LoginPrefix login action prefix. For example,"/login" in https://www.example.com/auth/login.
	LoginPrefix string
	//AuthPrefix auth action prefix. For example,"/auth" in https://www.example.com/auth/auth.
	AuthPrefix string
	//NotFoundAction action executed when provider keyword not found or auth params error.
	NotFoundAction func(w http.ResponseWriter, r *http.Request)
	//Session session used by auth provider to store data like state.
	Session Session
}

//Init auth with given option.
func (a *Auth) Init(option Option) error {
	return option.ApplyTo(a)
}

//ProviderManager return auth service's provider manager.
//Return MapProviderManager by default.
func (a *Auth) ProviderManager() ProviderManager {
	if a.providerManager != nil {
		return a.providerManager
	}
	return DefaultProviderManager
}

//SetProviderManager auth service's provider manager.
func (a *Auth) SetProviderManager(p ProviderManager) {
	a.providerManager = p
}

//RegisterProvider register provider with given keyword and driver.
//Return provider and any error if raised.
func (a *Auth) RegisterProvider(keyword string, driver Driver) (*Provider, error) {
	return a.ProviderManager().RegisterProvider(a, keyword, driver)
}

//RegisterWithCreator register provider with given keyword and driver creator.
//Return provider and any error if raised.
func (a *Auth) RegisterWithCreator(keyword string, creator DriverCreator) (*Provider, error) {
	return a.ProviderManager().RegisterProvider(a, keyword, creator.Create())
}

//MustRegisterProvider register provider with given keyword and driver.
//Panic if any error raised.
//Return provider.
func (a *Auth) MustRegisterProvider(keyword string, driver Driver) *Provider {
	s, err := a.ProviderManager().RegisterProvider(a, keyword, driver)
	if err != nil {
		panic(err)
	}
	return s
}

//MustRegisterWithCreator register provider with given keyword and driver creator.
//Panic if any error raised.
//Return provider.
func (a *Auth) MustRegisterWithCreator(keyword string, creator DriverCreator) *Provider {
	s, err := a.ProviderManager().RegisterProvider(a, keyword, creator.Create())
	if err != nil {
		panic(err)
	}
	return s
}

//GetProvider get registered provider by keyword.
//return provider and any error if raised.
func (a *Auth) GetProvider(keyword string) (*Provider, error) {
	return a.ProviderManager().GetProvider(a, keyword)
}

//MustGetProvider get registered provider by keyword.
//Panic if any error raised.
//return provider.
func (a *Auth) MustGetProvider(keyword string) *Provider {
	s, err := a.ProviderManager().GetProvider(a, keyword)
	if err != nil {
		panic(err)
	}
	return s
}

//New create new auth service.
//You should init auth after create.
func New() *Auth {
	return &Auth{}
}

//MustGetResult get result form http request.Create result if request result does not exist.
//Return result.
func (a *Auth) MustGetResult(req *http.Request) *Result {
	data := req.Context().Value(ResultContextName)
	if data != nil {
		result, ok := data.(*Result)
		if ok {
			return result
		}
	}
	return NewResult()
}

//SetResult set result into http request.
func (a *Auth) SetResult(r *http.Request, result *Result) {
	ctx := context.WithValue(r.Context(), ResultContextName, result)
	*r = *r.WithContext(ctx)
}

//RandToken generate random bytes in give length.
//Return random butes and any error if raised.
func (a *Auth) RandToken(length int) ([]byte, error) {
	token := make([]byte, length)
	_, err := rand.Read(token)
	if err != nil {
		return nil, err
	}
	l := len(TokenMask)
	for k, v := range token {
		index := int(v) % l
		token[k] = TokenMask[index]
	}
	return token, nil
}

//Serve serve as http handlerFunc.
//For example http.StripPrefix("/auth", http.HandlerFunc(Auth.Serve(successAction)))
//Params SuccessAction action execute after auth success.You should get auth result by Auth.MustGetResult method.
//Return http.handleFunc
func (a *Auth) Serve(SuccessAction func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		var provider *Provider
		var keyword string
		path := r.URL.Path
		if keyword = strings.TrimPrefix(path, a.LoginPrefix); len(keyword) < len(path) {
			provider, err = a.GetProvider(keyword)
			if err != nil {
				panic(err)
			}
			if provider != nil {
				provider.Login(w, r)
				return
			}
		} else if keyword = strings.TrimPrefix(path, a.AuthPrefix); len(keyword) < len(path) {
			provider, err = a.GetProvider(keyword)
			if err != nil {
				panic(err)
			}
			if provider != nil {
				result, err := provider.AuthRequest(r)
				if err == ErrAuthParamsError {
					if a.Debug {
						fmt.Println(ErrAuthParamsError.Error())
					}
					a.NotFoundAction(w, r)
					return
				}
				if err != nil {
					panic(err)
				}
				if result != nil && result.Account != "" {
					result.Keyword = provider.Keyword
					a.SetResult(r, result)
					SuccessAction(w, r)
					return
				}
				if a.Debug {
					fmt.Println("auth:user not found")
				}
			}
		}
		a.NotFoundAction(w, r)
	}
}
