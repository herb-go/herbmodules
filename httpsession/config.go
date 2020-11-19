package httpsession

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

func (c *Config) ApplyTo(s *Store) error {
	s.AutoStart = c.AutoStart
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
