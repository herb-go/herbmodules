package httpinfomanager

import (
	"github.com/herb-go/herb/middleware/httpinfo"
)

type ExtractorConfig struct {
	Name   string
	Type   string
	Config func(v interface{}) error `config:", lazyload"`
}

func (c *ExtractorConfig) Register() error {
	f, err := GetExtractorFactory(c.Type)
	if err != nil {
		return err
	}
	e, err := f.CreateExtractor(c.Config)
	if err != nil {
		return err
	}
	RegisterExtractor(c.Name, e)
	return nil
}

type FormatterConfig struct {
	Name   string
	Type   string
	Config func(v interface{}) error `config:", lazyload"`
}

func (c *FormatterConfig) Register() error {
	f, err := GetFormatterFactory(c.Type)
	if err != nil {
		return err
	}
	formatter, err := f.CreateFormatter(c.Config)
	if err != nil {
		return err
	}
	RegisterFormatter(c.Name, formatter)
	return nil
}

type FieldConfig struct {
	Name       FieldName
	Extractor  string
	Formatters []string
}

func (c *FieldConfig) ApplyTo() error {
	var err error
	if c.Name == "" {
		return ErrEmptyFieldName
	}
	e, err := GetExtractor(c.Extractor)
	if err != nil {
		return err
	}
	formatters := make([]httpinfo.Formatter, len(c.Formatters))
	for k, v := range c.Formatters {
		formatters[k], err = GetFormatter(v)
		if err != nil {
			return err
		}
	}
	field := httpinfo.NewExtractorField()
	field.Extractor = e
	field.Formatters = formatters
	RegisterField(string(c.Name), field)
	return nil
}

type FieldName string

func (n FieldName) Field() (httpinfo.Field, error) {
	return GetField(string(n))
}

type Config struct {
	Extractors []*ExtractorConfig
	Formatters []*FormatterConfig
	Fields     []*FieldConfig
}

func (c *Config) Register() error {
	for _, v := range c.Extractors {
		err := v.Register()
		if err != nil {
			return err
		}
	}
	for _, v := range c.Formatters {
		err := v.Register()
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Config) ApplyToFields() error {
	for _, v := range c.Fields {
		err := v.ApplyTo()
		if err != nil {
			return err
		}
	}
	return nil
}
