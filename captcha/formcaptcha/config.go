package formcaptcha

type Config struct {
	Name   string
	Config func(v interface{}) error `config:", lazyload"`
}

func (c *Config) ApplyTo(fc *FormCaptcha) error {
	fc.Name = c.Name
	d, err := NewDriver(c.Name, c.Config)
	if err != nil {
		return err
	}
	fc.Driver = d
	return nil
}
