package herbsession

type EngineName string
type Engine interface {
	EngineName() EngineName
	NewToken() (token string, err error)
	SessionTimeout() int64
	TokenLastActive(token string) (int64, error)
	LoadToken(token string) (newtoken string, data []byte, err error)
	SaveToken(token string, data []byte, maxLiftimeInSecond int64) (newtoken string, err error)
	RevokeToken(token string) (err error)
	DynamicToken() bool
}
