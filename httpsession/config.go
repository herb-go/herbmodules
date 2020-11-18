package httpsession

type Config struct {
	AutoStart   bool
	Engine      EngineConfig
	MaxLifetime int64
	Installer   *InstallerConfig
}

func (c *Config) ApplyTo(s *Store) error {
	s.AutoStart = c.AutoStart
	s.MaxLifetime = c.MaxLifetime
	e, err := c.Engine.CreateEngine()
	if err != nil {
		return err
	}
	s.Engine = e
	if c.Installer != nil {
		i, err := c.Installer.CreateInstaller()
		if err != nil {
			return err
		}
		s.Installer = i
	}
	return nil
}

type InstallerConfig struct {
	Name   InstallerName
	Config func(v interface{}) error `config:", lazyload"`
}

func (c *InstallerConfig) CreateInstaller() (Installer, error) {
	return CreateInstaller(c.Name, c.Config)
}

type EngineConfig struct {
	Name   EngineName
	Config func(v interface{}) error `config:", lazyload"`
}

func (c *EngineConfig) CreateEngine() (Engine, error) {
	return CreateEngine(c.Name, c.Config)
}
