package httpinfomanager_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/herb-go/herb/middleware/httpinfo"
	"github.com/herb-go/herbmodules/httpinfomanager"
)

var errTest = errors.New("test")

var testDefaultFormatterFactory = func(string) (httpinfomanager.FormatterFactory, error) {
	return nil, errTest
}
var testDefaultExtractorFactory = func(string) (httpinfomanager.ExtractorFactory, error) {
	return nil, errTest
}
var testDefaultValidatorFactory = func(string) (httpinfomanager.ValidatorFactory, error) {
	return nil, errTest
}

var testValidator = httpinfo.ValidatorAlways

var testValidatorFactory = httpinfomanager.ValidatorFactoryFunc(func(loader func(interface{}) error) (httpinfo.Validator, error) {
	return testValidator, nil
})
var testFormatter = httpinfo.FormatterFunc(func([]byte) ([]byte, bool, error) {
	return []byte("test"), true, nil
})

var testFormatterFactory = httpinfomanager.FormatterFactoryFunc(func(loader func(interface{}) error) (httpinfo.Formatter, error) {
	return testFormatter, nil
})

var testExtractor = httpinfo.ExtractorFunc(func(r *http.Request) ([]byte, error) {
	return []byte("test"), nil
})

var testExtractorFactory = httpinfomanager.ExtractorFactoryFunc(func(loader func(interface{}) error) (httpinfo.Extractor, error) {
	return testExtractor, nil
})

func TestRegister(t *testing.T) {
	httpinfomanager.Reset()
	defer httpinfomanager.Reset()
	ef, err := httpinfomanager.GetExtractorFactory("notexist")
	if ef != nil || err != httpinfomanager.ErrExtractorFactoryNotFound {
		t.Fatal(ef, err)
	}
	httpinfomanager.SetDefaultExtractorFactory(testDefaultExtractorFactory)
	ef, err = httpinfomanager.GetExtractorFactory("notexist")
	if ef != nil || err != errTest {
		t.Fatal(ef, err)
	}
	extractor, err := httpinfomanager.GetExtractor("notexist")
	if extractor != nil || errors.Unwrap(err) != httpinfomanager.ErrExtractorNotFound {
		t.Fatal(extractor, err)
	}
	httpinfomanager.RegisterExtractor("test", testExtractor)
	extractor, err = httpinfomanager.GetExtractor("test")
	if extractor == nil || err != nil {
		t.Fatal(extractor, err)
	}
	httpinfomanager.RegisterExtractorFactory("testfactory", testExtractorFactory)
	ef, err = httpinfomanager.GetExtractorFactory("testfactory")
	if ef == nil || err != nil {
		t.Fatal(ef, err)
	}
	extractor, err = ef.CreateExtractor(nil)
	if extractor == nil || err != nil {
		t.Fatal(extractor, err)
	}

	ff, err := httpinfomanager.GetFormatterFactory("notexist")
	if ff != nil || err != httpinfomanager.ErrFormatterFactoryNotFound {
		t.Fatal(ff, err)
	}
	httpinfomanager.SetDefaultFormatterFactory(testDefaultFormatterFactory)
	ff, err = httpinfomanager.GetFormatterFactory("notexist")
	if ff != nil || err != errTest {
		t.Fatal(ff, err)
	}
	formatter, err := httpinfomanager.GetFormatter("notexist")
	if formatter != nil || errors.Unwrap(err) != httpinfomanager.ErrFormatterNotFound {
		t.Fatal(formatter, err)
	}

	httpinfomanager.RegisterFormatter("test", testFormatter)
	formatter, err = httpinfomanager.GetFormatter("test")
	if formatter == nil || err != nil {
		t.Fatal(formatter, err)
	}
	httpinfomanager.RegisterFormatterFactory("testfactory", testFormatterFactory)
	ff, err = httpinfomanager.GetFormatterFactory("testfactory")
	if ff == nil || err != nil {
		t.Fatal(ff, err)
	}
	formatter, err = ff.CreateFormatter(nil)
	if formatter == nil || err != nil {
		t.Fatal(formatter, err)
	}

	vf, err := httpinfomanager.GetValidatorFactory("notexist")
	if vf != nil || err != httpinfomanager.ErrValidatorFactoryNotFound {
		t.Fatal(vf, err)
	}
	httpinfomanager.SetDefaultValidatorFactory(testDefaultValidatorFactory)
	vf, err = httpinfomanager.GetValidatorFactory("notexist")
	if vf != nil || err != errTest {
		t.Fatal(vf, err)
	}
	validator, err := httpinfomanager.GetValidator("notexist")
	if validator != nil || errors.Unwrap(err) != httpinfomanager.ErrValidatorNotFound {
		t.Fatal(validator, err)
	}
	httpinfomanager.RegisterValidator("test", testValidator)
	validator, err = httpinfomanager.GetValidator("test")
	if validator == nil || err != nil {
		t.Fatal(validator, err)
	}
	httpinfomanager.RegisterValidatorFactory("testfactory", testValidatorFactory)
	vf, err = httpinfomanager.GetValidatorFactory("testfactory")
	if vf == nil || err != nil {
		t.Fatal(vf, err)
	}
	validator, err = vf.CreateValidator(nil)
	if validator == nil || err != nil {
		t.Fatal(validator, err)
	}
}
