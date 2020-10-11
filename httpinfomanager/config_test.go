package httpinfomanager_test

import (
	"errors"
	"testing"

	"github.com/herb-go/herb/middleware/httpinfo"

	"github.com/herb-go/herbmodules/httpinfomanager"
)

var errFormatterFactory = httpinfomanager.FormatterFactoryFunc(func(loader func(interface{}) error) (httpinfo.Formatter, error) {
	return nil, errTest
})

var errExtractorFactory = httpinfomanager.ExtractorFactoryFunc(func(loader func(interface{}) error) (httpinfo.Extractor, error) {
	return nil, errTest
})

func TestField(t *testing.T) {
	httpinfomanager.Reset()
	defer httpinfomanager.Reset()
	fn := httpinfomanager.FieldName("not exists")
	field, err := fn.Field()
	if field != nil || errors.Unwrap(err) != httpinfomanager.ErrFieldNotFound {
		t.Fatal(field, err)
	}
}
func TestConfig(t *testing.T) {
	var err error
	httpinfomanager.Reset()
	defer httpinfomanager.Reset()
	ec := &httpinfomanager.ExtractorConfig{
		Name: "test",
		Type: "test",
	}
	err = ec.Register()
	if err != httpinfomanager.ErrExtractorFactoryNotFound {
		t.Fatal(err)
	}
	httpinfomanager.RegisterExtractorFactory("test", testExtractorFactory)
	err = ec.Register()
	if err != nil {
		t.Fatal(err)
	}
	fc := &httpinfomanager.FormatterConfig{
		Name: "test",
		Type: "test",
	}
	err = fc.Register()
	if err != httpinfomanager.ErrFormatterFactoryNotFound {
		t.Fatal(err)
	}
	httpinfomanager.RegisterFormatterFactory("test", testFormatterFactory)
	err = fc.Register()
	if err != nil {
		t.Fatal(err)
	}
	fieldconfig := &httpinfomanager.FieldConfig{}
	err = fieldconfig.ApplyTo()
	if err != httpinfomanager.ErrEmptyFieldName {
		t.Fatal(err)
	}
	fieldconfig = &httpinfomanager.FieldConfig{
		Name:      "notexist",
		Extractor: "notexist",
	}
	err = fieldconfig.ApplyTo()
	if errors.Unwrap(err) != httpinfomanager.ErrExtractorNotFound {
		t.Fatal(err)
	}
	fieldconfig = &httpinfomanager.FieldConfig{
		Name:       "notexist",
		Extractor:  "test",
		Formatters: []string{"notexist"},
	}
	err = fieldconfig.ApplyTo()
	if errors.Unwrap(err) != httpinfomanager.ErrFormatterNotFound {
		t.Fatal(err)
	}
	fieldconfig = &httpinfomanager.FieldConfig{
		Name:       "test",
		Extractor:  "test",
		Formatters: []string{"test"},
	}
	err = fieldconfig.ApplyTo()
	if err != nil {
		t.Fatal(err)
	}
	fn := httpinfomanager.FieldName("test")
	field, err := fn.Field()
	if err != nil {
		t.Fatal(err)
	}
	data, ok, err := field.LoadInfo(nil)
	if string(data) != "test" || ok != true || err != nil {
		t.Fatal(field)
	}
	httpinfomanager.RegisterExtractorFactory("err", errExtractorFactory)
	ec = &httpinfomanager.ExtractorConfig{
		Name: "err",
		Type: "err",
	}
	err = ec.Register()
	if err != errTest {
		t.Fatal(err)
	}
	httpinfomanager.RegisterFormatterFactory("err", errFormatterFactory)
	fc = &httpinfomanager.FormatterConfig{
		Name: "err",
		Type: "err",
	}
	err = fc.Register()
	if err != errTest {
		t.Fatal(err)
	}
}

func TestConfigs(t *testing.T) {
	httpinfomanager.Reset()
	defer httpinfomanager.Reset()
	ec := &httpinfomanager.ExtractorConfig{
		Name: "test",
		Type: "test",
	}
	err := ec.Register()
	if err != httpinfomanager.ErrExtractorFactoryNotFound {
		t.Fatal(err)
	}
	httpinfomanager.RegisterExtractorFactory("test", testExtractorFactory)
	err = ec.Register()
	if err != nil {
		t.Fatal(err)
	}
	fc := &httpinfomanager.FormatterConfig{
		Name: "test",
		Type: "test",
	}
	err = fc.Register()
	if err != httpinfomanager.ErrFormatterFactoryNotFound {
		t.Fatal(err)
	}
	httpinfomanager.RegisterFormatterFactory("test", testFormatterFactory)
	err = fc.Register()
	if err != nil {
		t.Fatal(err)
	}
	presetConfig := &httpinfomanager.Config{
		Extractors: []*httpinfomanager.ExtractorConfig{
			&httpinfomanager.ExtractorConfig{
				Name: "test",
				Type: "test",
			},
		},
		Formatters: []*httpinfomanager.FormatterConfig{
			&httpinfomanager.FormatterConfig{
				Name: "test",
				Type: "test",
			},
		},
	}
	config := &httpinfomanager.Config{
		Fields: []*httpinfomanager.FieldConfig{
			&httpinfomanager.FieldConfig{
				Name:       "test",
				Extractor:  "test",
				Formatters: []string{"test"},
			},
		},
	}
	err = presetConfig.Register()
	if err != nil {
		panic(err)
	}
	err = config.Register()
	if err != nil {
		panic(err)
	}
	err = presetConfig.ApplyToFields()
	if err != nil {
		panic(err)
	}
	err = config.ApplyToFields()
	if err != nil {
		panic(err)
	}
	fn := httpinfomanager.FieldName("test")
	field, err := fn.Field()
	if err != nil {
		panic(err)
	}
	data, ok, err := field.LoadInfo(nil)
	if string(data) != "test" || ok != true || err != nil {
		t.Fatal(field)
	}
}

func TestErrorConfig(t *testing.T) {
	httpinfomanager.Reset()
	defer httpinfomanager.Reset()
	var config *httpinfomanager.Config
	var err error
	config = &httpinfomanager.Config{
		Extractors: []*httpinfomanager.ExtractorConfig{
			&httpinfomanager.ExtractorConfig{
				Name: "test",
				Type: "test",
			},
		},
	}
	err = config.Register()
	if err == nil {
		t.Fatal(err)
	}
	config = &httpinfomanager.Config{
		Formatters: []*httpinfomanager.FormatterConfig{
			&httpinfomanager.FormatterConfig{
				Name: "test",
				Type: "test",
			},
		},
	}
	err = config.Register()
	if err == nil {
		t.Fatal(err)
	}
	config = &httpinfomanager.Config{
		Fields: []*httpinfomanager.FieldConfig{
			&httpinfomanager.FieldConfig{
				Name: "",
			},
		},
	}
	err = config.ApplyToFields()
	if err == nil {
		t.Fatal(err)
	}
}
