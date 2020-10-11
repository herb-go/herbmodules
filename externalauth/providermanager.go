package auth

//ProviderManager provider manager interface.
type ProviderManager interface {
	//GetProvider get provider by keyword and auth service
	//return provider and any erros if raised.
	GetProvider(auth *Auth, keyword string) (*Provider, error)
	//RegisterProvider register provider with keyword ,auth service and driver.
	//return provider and any erros if raised.
	RegisterProvider(auth *Auth, keyword string, driver Driver) (*Provider, error)
}

//MapProviderManager provider manager which store providers in map
type MapProviderManager struct {
	Providers map[string]*Provider
}

//NewMapProviderManager create new MapProviderManager
func NewMapProviderManager() *MapProviderManager {
	return &MapProviderManager{
		Providers: map[string]*Provider{},
	}
}

//GetProvider get provider by keyword and auth service
//return provider and any erros if raised.
func (m *MapProviderManager) GetProvider(a *Auth, keyword string) (*Provider, error) {
	s, ok := m.Providers[keyword]
	if ok {
		return s, nil
	}
	return nil, nil

}

//RegisterProvider register provider with keyword ,auth service and driver.
//return provider and any erros if raised.
func (m *MapProviderManager) RegisterProvider(a *Auth, keyword string, driver Driver) (*Provider, error) {
	s := &Provider{
		Driver:  driver,
		Auth:    a,
		Keyword: keyword,
	}
	m.Providers[keyword] = s
	return s, nil

}

//DefaultProviderManager default provider manager for auth service.
var DefaultProviderManager = NewMapProviderManager()
