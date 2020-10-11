package auth

type testProviderManager struct {
	provider *Provider
}

func (m *testProviderManager) GetProvider(auth *Auth, keyword string) (*Provider, error) {
	return m.provider, nil
}
func (m *testProviderManager) RegisterProvider(auth *Auth, keyword string, driver Driver) (*Provider, error) {
	m.provider = &Provider{
		Driver:  driver,
		Auth:    auth,
		Keyword: keyword,
	}
	return m.provider, nil
}

func newTestProviderManager() *testProviderManager {
	return &testProviderManager{}
}
