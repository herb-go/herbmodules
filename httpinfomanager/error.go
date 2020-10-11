package httpinfomanager

import "errors"

var ErrExtractorNotFound = errors.New("extractor not found")

var ErrExtractorFactoryNotFound = errors.New("extractor factory not found")

var ErrFormatterNotFound = errors.New("formatter not found")

var ErrFormatterFactoryNotFound = errors.New("formatter factory not found")

var ErrValidatorNotFound = errors.New("validator not found")

var ErrFieldNotFound = errors.New("field not found")

var ErrEmptyFieldName = errors.New("empty field name")

var ErrValidatorFactoryNotFound = errors.New("validator factory not found")

var ErrUnavailableIndex = errors.New("unavailable index")
