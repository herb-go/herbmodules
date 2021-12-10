package httpsession

var DefaultTimeout = int64(15 * 60)
var DefaultMaxLifetime = int64(7 * 24 * 3600)
var DefaultLastActiveInterval = int64(60)

type Config struct {
	AutoStart          bool
	Engine             EngineName
	EngineConfig       func(interface{}) error `config:", lazyload"`
	MaxLifetime        int64
	Timeout            int64
	LastActiveInterval int64
	Installer          InstallerName
	InstallerConfig    func(interface{}) error `config:", lazyload"`
}

func setConfigDefaults(c *Config) {
	if c.MaxLifetime == 0 {
		c.MaxLifetime = DefaultMaxLifetime
	}
	if c.Timeout == 0 {
		c.Timeout = DefaultTimeout
	}
	if c.LastActiveInterval == 0 {
		c.LastActiveInterval = DefaultLastActiveInterval
	}
}
func (c *Config) ApplyTo(s *Store) error {
	setConfigDefaults(c)
	s.AutoStart = !c.AutoStart
	s.MaxLifetime = c.MaxLifetime
	s.Timeout = c.Timeout
	s.LastActiveInterval = c.LastActiveInterval
	e, err := CreateEngine(c.Engine, c.EngineConfig)
	if err != nil {
		return err
	}
	s.Engine = e
	if c.Installer != "" {
		i, err := CreateInstaller(c.Installer, c.InstallerConfig)
		if err != nil {
			return err
		}
		s.Installer = i
	}
	return nil
}
