package httpinfomanager

import (
	"fmt"
	"sync"

	"github.com/herb-go/herb/middleware/httpinfo"
)

var lock sync.Mutex
var extractorFactories = map[string]ExtractorFactory{}
var defaultExtractorFactory func(string) (ExtractorFactory, error)

func notFoundDefaultExtractorFactory(string) (ExtractorFactory, error) {
	return nil, ErrExtractorFactoryNotFound
}

var formatterFactories = map[string]FormatterFactory{}
var defaultFormatterFactory func(string) (FormatterFactory, error) = notFoundDefaultFormatterFactory

func notFoundDefaultFormatterFactory(string) (FormatterFactory, error) {
	return nil, ErrFormatterFactoryNotFound
}

var validatorFactories = map[string]ValidatorFactory{}

var defaultValidatorFactory func(string) (ValidatorFactory, error) = notFoundDefaultValidatorFactory

func notFoundDefaultValidatorFactory(string) (ValidatorFactory, error) {
	return nil, ErrValidatorFactoryNotFound
}

var registeredExtractors = map[string]httpinfo.Extractor{}
var registeredFormatters = map[string]httpinfo.Formatter{}
var registeredValidators = map[string]httpinfo.Validator{}

var registeredFields = map[string]httpinfo.Field{}

func Reset() {
	lock.Lock()
	defer lock.Unlock()
	extractorFactories = map[string]ExtractorFactory{}
	formatterFactories = map[string]FormatterFactory{}
	validatorFactories = map[string]ValidatorFactory{}
	registeredExtractors = map[string]httpinfo.Extractor{}
	registeredFormatters = map[string]httpinfo.Formatter{}
	registeredValidators = map[string]httpinfo.Validator{}
	defaultFormatterFactory = notFoundDefaultFormatterFactory
	defaultExtractorFactory = notFoundDefaultExtractorFactory
	defaultValidatorFactory = notFoundDefaultValidatorFactory
	registeredFields = map[string]httpinfo.Field{}
}
func SetDefaultExtractorFactory(f func(string) (ExtractorFactory, error)) {
	lock.Lock()
	defer lock.Unlock()
	defaultExtractorFactory = f
}
func SetDefaultFormatterFactory(f func(string) (FormatterFactory, error)) {
	lock.Lock()
	defer lock.Unlock()
	defaultFormatterFactory = f
}

func SetDefaultValidatorFactory(f func(string) (ValidatorFactory, error)) {
	lock.Lock()
	defer lock.Unlock()
	defaultValidatorFactory = f
}

func RegisterExtractor(name string, e httpinfo.Extractor) {
	lock.Lock()
	defer lock.Unlock()
	registeredExtractors[name] = e
}
func GetExtractor(name string) (httpinfo.Extractor, error) {
	lock.Lock()
	defer lock.Unlock()
	r := registeredExtractors[name]
	if r == nil {
		return nil, fmt.Errorf("httpinfomanager: %w (%s)", ErrExtractorNotFound, name)
	}
	return r, nil
}

func GetExtractorFactory(name string) (ExtractorFactory, error) {
	lock.Lock()
	defer lock.Unlock()
	f, ok := extractorFactories[name]
	if ok {
		return f, nil
	}
	return defaultExtractorFactory(name)
}
func RegisterExtractorFactory(name string, f ExtractorFactory) {
	lock.Lock()
	defer lock.Unlock()
	extractorFactories[name] = f
}
func RegisterFormatter(name string, formatter httpinfo.Formatter) {
	lock.Lock()
	defer lock.Unlock()
	registeredFormatters[name] = formatter
}

func GetFormatter(name string) (httpinfo.Formatter, error) {
	lock.Lock()
	defer lock.Unlock()
	f := registeredFormatters[name]
	if f == nil {
		return nil, fmt.Errorf("httpinfomanager: %w (%s)", ErrFormatterNotFound, name)
	}
	return f, nil
}

func GetFormatterFactory(name string) (FormatterFactory, error) {
	lock.Lock()
	defer lock.Unlock()
	f, ok := formatterFactories[name]
	if ok {
		return f, nil
	}
	return defaultFormatterFactory(name)
}
func RegisterFormatterFactory(name string, f FormatterFactory) {
	lock.Lock()
	defer lock.Unlock()
	formatterFactories[name] = f
}
func RegisterValidator(name string, validator httpinfo.Validator) {
	lock.Lock()
	defer lock.Unlock()
	registeredValidators[name] = validator
}

func GetValidator(name string) (httpinfo.Validator, error) {
	lock.Lock()
	defer lock.Unlock()
	f := registeredValidators[name]
	if f == nil {
		return nil, fmt.Errorf("httpinfomanager: %w (%s)", ErrValidatorNotFound, name)
	}
	return f, nil
}

func GetValidatorFactory(name string) (ValidatorFactory, error) {
	lock.Lock()
	defer lock.Unlock()

	f, ok := validatorFactories[name]
	if ok {
		return f, nil
	}
	return defaultValidatorFactory(name)
}

func RegisterValidatorFactory(name string, f ValidatorFactory) {
	lock.Lock()
	defer lock.Unlock()
	validatorFactories[name] = f
}

func GetField(name string) (httpinfo.Field, error) {
	lock.Lock()
	defer lock.Unlock()
	f, ok := registeredFields[name]
	if !ok {
		return nil, fmt.Errorf("httpinfomanager: %w (%s)", ErrFieldNotFound, name)
	}
	return f, nil
}
func MustGetField(name string) httpinfo.Field {
	f, err := GetField(name)
	if err != nil {
		panic(err)
	}
	return f
}
func GetStringField(name string) (*httpinfo.StringField, error) {
	f, err := GetField(name)
	if err != nil {
		return nil, err
	}
	return httpinfo.NewStringField(f), nil
}

func MustGetStringField(name string) *httpinfo.StringField {
	f, err := GetStringField(name)
	if err != nil {
		panic(err)
	}
	return f
}
func RegisterField(name string, field httpinfo.Field) {
	lock.Lock()
	defer lock.Unlock()
	registeredFields[name] = field
}
