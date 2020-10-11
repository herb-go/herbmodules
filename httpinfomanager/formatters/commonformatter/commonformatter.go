package commonformatter

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/herb-go/herb/middleware/httpinfo"
	"github.com/herb-go/herbmodules/httpinfomanager"
)

type StringsFormatter func(string) string

func (f StringsFormatter) Format(data []byte) ([]byte, bool, error) {
	return []byte(f(string(data))), true, nil
}

var ToUpperFormatter = StringsFormatter(strings.ToUpper)

var ToUpperFormatterFactory = httpinfomanager.FormatterFactoryFunc(func(loader func(interface{}) error) (httpinfo.Formatter, error) {
	return ToUpperFormatter, nil
})

var ToLowerFormatter = StringsFormatter(strings.ToLower)

var ToLowerFormatterFactory = httpinfomanager.FormatterFactoryFunc(func(loader func(interface{}) error) (httpinfo.Formatter, error) {
	return ToLowerFormatter, nil
})

var TrimSpaceFormatter = StringsFormatter(strings.TrimSpace)

var TrimSpaceFormatterFactory = httpinfomanager.FormatterFactoryFunc(func(loader func(interface{}) error) (httpinfo.Formatter, error) {
	return TrimSpaceFormatter, nil
})
var IntegerFormatter = httpinfo.FormatterFunc(func(data []byte) ([]byte, bool, error) {
	_, err := strconv.Atoi(string(data))
	if err != nil {
		return nil, false, nil
	}
	return data, true, nil
})
var IntegerFormatterFactory = httpinfomanager.FormatterFactoryFunc(func(loader func(interface{}) error) (httpinfo.Formatter, error) {
	return IntegerFormatter, nil
})

type SplitFormatter struct {
	Sep   string
	Index int
}

func (f *SplitFormatter) Format(data []byte) ([]byte, bool, error) {
	if f.Index < 0 {
		return nil, false, fmt.Errorf("httpinfomanager : %w (%d)", httpinfomanager.ErrUnavailableIndex, f.Index)
	}
	if f.Sep == "" {
		return nil, false, ErrEmptySep
	}
	results := strings.Split(string(data), f.Sep)
	if len(results) > f.Index {
		return []byte(results[f.Index]), true, nil
	}
	return nil, false, nil
}

var SplitFormatterFactory = httpinfomanager.FormatterFactoryFunc(func(loader func(interface{}) error) (httpinfo.Formatter, error) {
	f := &SplitFormatter{}
	err := loader(f)
	if err != nil {
		return nil, err
	}

	return f, nil
})

type regexpFormatter struct {
	Regexp *regexp.Regexp
	Index  int
}

func (f *regexpFormatter) match(data []byte) ([]byte, bool, error) {
	ok := f.Regexp.Match(data)
	if ok {
		return data, true, nil
	}
	return nil, false, nil
}

func (f *regexpFormatter) find(data []byte) ([]byte, bool, error) {
	results := f.Regexp.FindSubmatch(data)
	if len(results) > f.Index {
		return results[f.Index+1], true, nil
	}
	return nil, false, nil
}

type MatchFormatter struct {
	formatter *regexpFormatter
}

func (f MatchFormatter) Format(data []byte) ([]byte, bool, error) {
	return f.formatter.match(data)
}

type FindFormatter struct {
	formatter *regexpFormatter
}

func (f FindFormatter) Format(data []byte) ([]byte, bool, error) {
	return f.formatter.find(data)
}

type RegexpConfig struct {
	Pattern string
	Index   int
}

func (c *RegexpConfig) createFormatter() (*regexpFormatter, error) {
	if c.Index < 0 {
		return nil, fmt.Errorf("httpinfomanager : %w (%d)", httpinfomanager.ErrUnavailableIndex, c.Index)
	}
	p, err := regexp.Compile(c.Pattern)
	if err != nil {
		return nil, err
	}

	f := &regexpFormatter{
		Regexp: p,
		Index:  c.Index,
	}
	return f, nil
}

func (c *RegexpConfig) CreateMatchFormatter() (*MatchFormatter, error) {
	f, err := c.createFormatter()
	if err != nil {
		return nil, err
	}
	return &MatchFormatter{
		formatter: f,
	}, nil
}

func (c *RegexpConfig) CreateFindFormatter() (*FindFormatter, error) {
	f, err := c.createFormatter()
	if err != nil {
		return nil, err
	}
	return &FindFormatter{
		formatter: f,
	}, nil
}

var MatchFormatterFactory = httpinfomanager.FormatterFactoryFunc(func(loader func(interface{}) error) (httpinfo.Formatter, error) {
	c := &RegexpConfig{}
	err := loader(c)
	if err != nil {
		return nil, err
	}
	return c.CreateMatchFormatter()
})

var FindFormatterFactory = httpinfomanager.FormatterFactoryFunc(func(loader func(interface{}) error) (httpinfo.Formatter, error) {
	c := &RegexpConfig{}
	err := loader(c)
	if err != nil {
		return nil, err
	}
	return c.CreateFindFormatter()
})

func RegisterFactories() {
	httpinfomanager.RegisterFormatterFactory("toupper", ToUpperFormatterFactory)
	httpinfomanager.RegisterFormatterFactory("tolower", ToLowerFormatterFactory)
	httpinfomanager.RegisterFormatterFactory("trim", TrimSpaceFormatterFactory)
	httpinfomanager.RegisterFormatterFactory("integer", IntegerFormatterFactory)
	httpinfomanager.RegisterFormatterFactory("match", MatchFormatterFactory)
	httpinfomanager.RegisterFormatterFactory("find", FindFormatterFactory)
	httpinfomanager.RegisterFormatterFactory("split", SplitFormatterFactory)
}
