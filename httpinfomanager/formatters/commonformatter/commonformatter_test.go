package commonformatter_test

import (
	"encoding/json"
	"testing"

	"github.com/herb-go/herbconfig/loader"
	_ "github.com/herb-go/herbconfig/loader/drivers/jsonconfig"
	"github.com/herb-go/herbmodules/httpinfomanager"
	"github.com/herb-go/herbmodules/httpinfomanager/formatters/commonformatter"
	"github.com/herb-go/worker"
)

func newLoader(v interface{}) func(interface{}) error {
	bytes, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return loader.NewLoader("json", bytes)
}
func format(name string, loader func(interface{}) error, data []byte) ([]byte, bool, error) {
	f, err := httpinfomanager.GetFormatterFactory(name)
	if err != nil {
		return nil, false, err
	}
	formatter, err := f.CreateFormatter(loader)
	if err != nil {
		return nil, false, err
	}
	return formatter.Format(data)
}
func formatString(name string, loader func(interface{}) error, datastr string) (string, bool, error) {
	data, ok, err := format(name, loader, []byte(datastr))
	return string(data), ok, err
}
func TestFormatter(t *testing.T) {
	// var data []byte
	var datastr string
	var ok bool
	var err error

	worker.Reset()
	defer worker.Reset()
	httpinfomanager.Reset()
	defer httpinfomanager.Reset()
	commonformatter.RegisterFactories()
	datastr, ok, err = formatString("test.notfound", nil, "abc")
	if datastr != "" || ok != false || err == nil {
		t.Fatal(datastr, ok, err)
	}
	datastr, ok, err = formatString("toupper", nil, "abc")
	if datastr != "ABC" || ok != true || err != nil {
		t.Fatal(datastr, ok, err)
	}
	datastr, ok, err = formatString("tolower", nil, "ABC")
	if datastr != "abc" || ok != true || err != nil {
		t.Fatal(datastr, ok, err)
	}
	datastr, ok, err = formatString("trim", nil, " abc ")
	if datastr != "abc" || ok != true || err != nil {
		t.Fatal(datastr, ok, err)
	}
	datastr, ok, err = formatString("integer", nil, " abc ")
	if datastr != "" || ok != false || err != nil {
		t.Fatal(datastr, ok, err)
	}
	datastr, ok, err = formatString("integer", nil, "12345")
	if datastr != "12345" || ok != true || err != nil {
		t.Fatal(datastr, ok, err)
	}
	regexpconfig := commonformatter.RegexpConfig{
		Pattern: "^(abc)(def)$",
		Index:   1,
	}
	datastr, ok, err = formatString("match", newLoader(regexpconfig), "12345")
	if datastr != "" || ok != false || err != nil {
		t.Fatal(datastr, ok, err)
	}
	datastr, ok, err = formatString("match", newLoader(regexpconfig), "abcdef")
	if datastr != "abcdef" || ok != true || err != nil {
		t.Fatal(datastr, ok, err)
	}
	datastr, ok, err = formatString("find", newLoader(regexpconfig), "12345")
	if datastr != "" || ok != false || err != nil {
		t.Fatal(datastr, ok, err)
	}
	datastr, ok, err = formatString("find", newLoader(regexpconfig), "abcdef")
	if datastr != "def" || ok != true || err != nil {
		t.Fatal(datastr, ok, err)
	}
	regexpwrongconfig := commonformatter.RegexpConfig{
		Pattern: "^(abc)(def)$",
		Index:   -1,
	}
	datastr, ok, err = formatString("find", newLoader(regexpwrongconfig), "12345")
	if datastr != "" || ok != false || err == nil {
		t.Fatal(datastr, ok, err)
	}
	regexpnotfounndconfig := commonformatter.RegexpConfig{
		Pattern: "^(abc)(def)$",
		Index:   2,
	}
	datastr, ok, err = formatString("find", newLoader(regexpnotfounndconfig), "12345")
	if datastr != "" || ok != false || err != nil {
		t.Fatal(datastr, ok, err)
	}
	splitconfig := &commonformatter.SplitFormatter{
		Sep:   "-",
		Index: 0,
	}
	datastr, ok, err = formatString("split", newLoader(splitconfig), "123-45")
	if datastr != "123" || ok != true || err != nil {
		t.Fatal(datastr, ok, err)
	}
	splitnotfoundconfig := &commonformatter.SplitFormatter{
		Sep:   "-",
		Index: 2,
	}
	datastr, ok, err = formatString("split", newLoader(splitnotfoundconfig), "123-45")
	if datastr != "" || ok != false || err != nil {
		t.Fatal(datastr, ok, err)
	}
	splitemptyconfig := &commonformatter.SplitFormatter{
		Index: 0,
	}
	datastr, ok, err = formatString("split", newLoader(splitemptyconfig), "123-45")
	if datastr != "" || ok != false || err == nil {
		t.Fatal(datastr, ok, err)
	}
	splitwrongconfig := &commonformatter.SplitFormatter{
		Sep:   "-",
		Index: -1,
	}
	datastr, ok, err = formatString("split", newLoader(splitwrongconfig), "123-45")
	if datastr != "" || ok != false || err == nil {
		t.Fatal(datastr, ok, err)
	}
}
