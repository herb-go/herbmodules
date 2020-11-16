package herbsession

type Engine interface {
	NewToken() (token string, err error)
	SessionTimeout() int64
	TokenLastActive(token string) (int64, error)
	LoadToken(token string) (newtoken string, data []byte, err error)
	UpdateToken(token string, data []byte, maxLiftimeInSecond int64) (newtoken string, err error)
	RevokeToken(token string) (newtoken string, err error)
	DynamicToken() bool
}
